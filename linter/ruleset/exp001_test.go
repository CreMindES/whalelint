package ruleset_test

import (
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/stretchr/testify/assert"

	RuleSet "github.com/cremindes/whalelint/linter/ruleset"
)

func TestValidateExp001(t *testing.T) {
	t.Parallel()

	// nolint:gofmt,gofumpt,goimports
	testCases := []struct {
		portValue    []string
		isViolation  bool
		name         string
	}{
		{portValue: []string{"4242"                        }, isViolation: false, name: "EXPOSE 4242"      },
		{portValue: []string{"4242/tcp"                    }, isViolation: false, name: "EXPOSE 4242/tcp"  },
		{portValue: []string{"4242/udp"                    }, isViolation: false, name: "EXPOSE 4242/udp"  },
		{portValue: []string{"4242/yyy"                    }, isViolation:  true, name: "EXPOSE 4242/yyy"  },
		{portValue: []string{"4242:tcp"                    }, isViolation:  true, name: "EXPOSE 4242:tcp"  },
		{portValue: []string{"4242", "4242/tcp", "4242/udp"}, isViolation: false,
			name: "EXPOSE 4242, 4242/tcp, 4242/udp"},
		{portValue: []string{"67999"                       }, isViolation:  true, name: "EXPOSE 67999"     },
		{portValue: []string{"4242", "67999", "4242/udp"   }, isViolation:  true,
			name: "EXPOSE 4242, 67999, 4242/udp"},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			command := &instructions.ExposeCommand{
				Ports: testCase.portValue,
			}

			assert.Equal(t, testCase.isViolation, RuleSet.ValidateExp001(command).IsViolated())
		})
	}
}
