package license

import (
	"fmt"
	"log"
	"testing"
)

func TestParse(t *testing.T) {
	content := "kKZBzS2H8V/v2/eaSfOS3FMORur/pLzgiz5nJ4grDdrYqPcpCvJaIeTaDfg8EaB4BjgKFd2q69CGQ2urzlmhtY8QNa2Gr324QcdC/XkzF49+FlBUDHsvsxpoOxLypC1RJpaqtdTP8tS3JkO4/Uf4seXNZ/Wa+NQXWnyzMspJ3g8kmZ8P1pMTUzt53i8g0cOXmpmT/s5J0fjl2AwINTxyHoXXWVAIgPf/oGjkl6rU3YdaWTWaIAsxbkwaas241vVl/eG/QnjWvfLnuubcy2SpBm+DIS2UYZpDsLmQYh9WpHKy5BvuU0X2HwwterMDASQoucDuGJcAkq1C0xXygLJ4B9AO+RyrioRxALlR9fvJa0Rr3d0tROuIVGhwajd1El94mJQ3n1LE+iIMuJfQGD6AmZQa0OT0tA6oUrXXtjgnX1OckL+zLYasCUG/opTycYljV4vLOdHwcj3p2n3VD6EQKk5idjV6ChX4IihiXLriw/tMlKm+9QQWMfUPvDLNPfWPPbRIehUv9eUh/zvQniC7Wc8Nd7f+oC5hd5KUXrmsz+yQpUysdGpB3MQakUSEEnYxt9S0UZ4qSIqvMGWY73tVPALJ38ieThEEQS9CXIrXITOqTN1t360U04/4uLha/WQz2JXV/SU6Y3pn2OJhu2JVlR8oZwyFtvzBUS5pvvK16Qy7KHA="
	resp, err := Parse(content)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)
}
