package validator

import (
	"fmt"
	"log/slog"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/bamboo-firewall/be/api/v1/dto"
	"github.com/bamboo-firewall/be/pkg/entity"
	"github.com/bamboo-firewall/be/pkg/httpbase"
	"github.com/bamboo-firewall/be/pkg/net"
	"github.com/bamboo-firewall/be/pkg/selector"
)

func registerValidator(tagName string, validatorFunc validator.Func) {
	if err := httpbase.RegisterValidator(tagName, validatorFunc); err != nil {
		panic(err)
	} else {
		slog.Debug("registered validator", "tag", tagName)
	}
}

func registerStructValidation(fn validator.StructLevelFunc, in ...interface{}) {
	httpbase.RegisterStructValidation(fn, in...)
}

func Init() {
	registerValidator("name", validateName)
	registerValidator("selector", validateSelector)
	registerValidator("action", validateAction)
	registerValidator("ip_version", validateIPVersion)
	registerValidator("protocol", validateProtocol)
	registerValidator("port", validatePort)

	registerValidator("net", validateIPNetwork)
	registerValidator("cidr", validateCIDR)
	registerValidator("ip", validateIP)

	registerStructValidation(validateGNPSpecInput, dto.GNPSpecInput{})
	registerStructValidation(validateGNPSpecRuleInput, dto.GNPSpecRuleInput{})
	registerStructValidation(validateGNPSpecRuleEntityInput, dto.GNPSpecRuleEntityInput{})
	registerStructValidation(validateGNSSpecInput, dto.GNSSpecInput{})
}

var nameRegex = regexp.MustCompile(`^[-a-zA-Z0-9_\\.]+$`)

func validateName(fl validator.FieldLevel) bool {
	return nameRegex.MatchString(fl.Field().String())
}

func validateSelector(fl validator.FieldLevel) bool {
	sel := fl.Field().Interface().(string)
	_, err := selector.Parse(sel)
	if err != nil {
		return false
	}
	return true
}

func validateAction(fl validator.FieldLevel) bool {
	action := fl.Field().Interface().(string)
	return slices.Contains(
		[]entity.RuleAction{entity.RuleActionAllow, entity.RuleActionDeny, entity.RuleActionLog, entity.RuleActionPass},
		entity.RuleAction(strings.ToLower(action)),
	)
}

func validateIPVersion(fl validator.FieldLevel) bool {
	ipVersion := fl.Field().Interface().(int)
	return slices.Contains([]int{entity.IPVersion4, entity.IPVersion6}, ipVersion)
}

func validateProtocol(fl validator.FieldLevel) bool {
	protocol := fl.Field().Interface()
	switch protocol.(type) {
	case string:
		return slices.Contains([]string{entity.ProtocolTCP, entity.ProtocolUDP, entity.ProtocolICMP, entity.ProtocolSCTP, entity.ProtocolUDPLite}, strings.ToLower(protocol.(string)))
	case float64:
		protocolNum := uint8(protocol.(float64))
		return protocolNum != 0
	default:
		return false
	}
}

func validateCIDR(fl validator.FieldLevel) bool {
	n := fl.Field().String()
	_, _, err := net.ParseCIDROrIP(n)
	return err == nil
}

func validateIPNetwork(fl validator.FieldLevel) bool {
	n := fl.Field().String()
	ip, ipnet, err := net.ParseCIDROrIP(n)
	if err != nil {
		return false
	}
	return ip.String() == ipnet.IP.String()
}

func validateIP(fl validator.FieldLevel) bool {
	return net.ParseIP(fl.Field().String()) != nil
}

var portRangeRegex = regexp.MustCompile(`^(\d+):(\d+)$`)

const (
	portRangeMin int = 0
	portRangeMax int = 65535
)

// validatePort port range 0-65535
func validatePort(fl validator.FieldLevel) bool {
	if portNumber, ok := fl.Field().Interface().(float64); ok {
		if int(portNumber) < portRangeMin || int(portNumber) > portRangeMax {
			return false
		}
	} else if portRange, ok := fl.Field().Interface().(string); ok {
		portsMatch := portRangeRegex.FindStringSubmatch(portRange)
		if portsMatch == nil {
			return false
		}
		portStart, err := strconv.ParseUint(portsMatch[1], 10, 16)
		if err != nil {
			return false
		}
		portEnd, err := strconv.ParseUint(portsMatch[2], 10, 16)
		if err != nil {
			return false
		}
		if portStart > portEnd {
			return false
		}
	} else {
		return false
	}
	return true
}

func validateGNPSpecInput(sl validator.StructLevel) {
	input := sl.Current().Interface().(dto.GNPSpecInput)
	if len(input.Ingress) == 0 && len(input.Egress) == 0 {
		sl.ReportError(input.Ingress, "ingress", "Egress", "require ingress or egress", "")
	}
}

func validateGNPSpecRuleInput(sl validator.StructLevel) {
	input := sl.Current().Interface().(dto.GNPSpecRuleInput)
	if input.Protocol != nil && input.NotProtocol != nil {
		sl.ReportError(input.NotProtocol, "notProtocol", "NotProtocol", "cannot use notProtocol with protocol", "")
	}
	if input.Protocol != nil || input.NotProtocol != nil {
		if (input.Protocol != nil && !isProtocolSupportPort(input.Protocol)) || (input.NotProtocol != nil && !isProtocolSupportPort(input.NotProtocol)) {
			if input.Source != nil {
				if len(input.Source.Ports) > 0 {
					sl.ReportError(input.Source.Ports, "notPorts", "NotPorts", "protocol not support ports", "")
				}
				if len(input.Source.NotPorts) > 0 {
					sl.ReportError(input.Source.NotPorts, "notPorts", "NotPorts", "protocol not support ports", "")
				}
			}

			if input.Destination != nil {
				if len(input.Destination.Ports) > 0 {
					sl.ReportError(input.Destination.Ports, "notPorts", "NotPorts", "protocol not support ports", "")
				}
				if len(input.Destination.NotPorts) > 0 {
					sl.ReportError(input.Destination.NotPorts, "notPorts", "NotPorts", "protocol not support ports", "")
				}
			}
		}
	}

	var (
		seenV4, seenV6 bool
	)
	var scanNets = func(nets []string, fieldName string) {
		var v4, v6 bool
		for i, ipNetwork := range nets {
			ip, ipnet, err := net.ParseCIDR(ipNetwork)
			if err != nil {
				sl.ReportError(ipNetwork, fmt.Sprintf("%s[%d]", fieldName, i), "", "net", "")
				continue
			}
			if ip.String() != ipnet.IP.String() {
				sl.ReportError(ipNetwork, fmt.Sprintf("%s[%d]", fieldName, i), "", "ip network is invalid", "")
			}
			if input.IPVersion != nil && ip.Version() != *input.IPVersion {
				sl.ReportError(ipNetwork, fmt.Sprintf("%s[%d]", fieldName, i), "", "not match with ipVersion", "")
			}

			v4 = v4 || ip.Version() == entity.IPVersion4
			v6 = v6 || ip.Version() == entity.IPVersion6
		}

		if v4 && seenV6 || v6 && seenV4 || v4 && v6 {
			sl.ReportError(nets, fieldName, "", "cannot use ipV4 and ipV6 together", "")
		}

		seenV4 = seenV4 || v4
		seenV6 = seenV6 || v6
	}

	if input.Source != nil {
		scanNets(input.Source.Nets, "Source.Nets")
		scanNets(input.Source.NotNets, "Source.NotNets")
	}
	if input.Destination != nil {
		scanNets(input.Destination.Nets, "Destination.Nets")
		scanNets(input.Destination.NotNets, "Destination.NotNets")
	}
}

func isProtocolSupportPort(protocol interface{}) bool {
	switch protocol.(type) {
	case string:
		return slices.Contains([]string{entity.ProtocolTCP, entity.ProtocolUDP, entity.ProtocolSCTP}, strings.ToLower(protocol.(string)))
	case float64:
		protocolNum := uint8(protocol.(float64))
		return protocolNum == entity.ProtocolNumTCP || protocolNum == entity.ProtocolNumUDP || protocolNum == entity.ProtocolNumSCTP
	default:
		return false
	}
}

func validateGNPSpecRuleEntityInput(sl validator.StructLevel) {
	input := sl.Current().Interface().(dto.GNPSpecRuleEntityInput)
	if len(input.Nets) > 0 && len(input.NotNets) > 0 {
		sl.ReportError(input.NotNets, "notNets", "NotNets", "cannot use notNets with nets", "")
	}
	if len(input.NotPorts) > 0 && len(input.NotPorts) > 0 {
		sl.ReportError(input.NotPorts, "notPorts", "NotPorts", "cannot use notPorts with ports", "")
	}
}

func validateGNSSpecInput(sl validator.StructLevel) {
	input := sl.Current().Interface().(dto.GNSSpecInput)
	cidrMap := make(map[string]struct{})
	for i, netString := range input.Nets {
		ip, ipnet, err := net.ParseCIDROrIP(netString)
		if err != nil {
			sl.ReportError(netString, fmt.Sprintf("nets[%d]", i), "", "cidr", "")
			continue
		}
		var netV4V6 string
		if ip.String() == ipnet.IP.String() {
			netV4V6 = ipnet.String()
		} else {
			netV4V6 = ip.Network().String()
		}
		if _, ok := cidrMap[netV4V6]; ok {
			sl.ReportError(netV4V6, fmt.Sprintf("nets[%d]", i), "", "duplicate", "")
		} else {
			cidrMap[netV4V6] = struct{}{}
		}
	}
}
