package kube

import (
	"testing"
	"time"

	"github.com/skupperproject/skupper/internal/cmd/skupper/common"
	"github.com/skupperproject/skupper/internal/cmd/skupper/common/testutils"
	"github.com/skupperproject/skupper/internal/cmd/skupper/common/utils"
	fakeclient "github.com/skupperproject/skupper/internal/kube/client/fake"
	"github.com/skupperproject/skupper/pkg/apis/skupper/v2alpha1"
	"gotest.tools/v3/assert"
	v12 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestCmdLinkUpdate_ValidateInput(t *testing.T) {
	type test struct {
		name           string
		args           []string
		flags          common.CommandLinkUpdateFlags
		k8sObjects     []runtime.Object
		skupperObjects []runtime.Object
		expectedError  string
		skupperError   string
	}

	testTable := []test{
		{
			name:          "missing CRD",
			args:          []string{"my-connector", "8080"},
			flags:         common.CommandLinkUpdateFlags{},
			skupperError:  utils.CrdErr,
			expectedError: utils.CrdHelpErr,
		},
		{
			name:  "link is not updated because there is no site in the namespace.",
			args:  []string{"my-link"},
			flags: common.CommandLinkUpdateFlags{Cost: "1", Timeout: time.Minute},
			skupperObjects: []runtime.Object{
				&v2alpha1.Link{
					ObjectMeta: v1.ObjectMeta{
						Name:      "my-link",
						Namespace: "test",
					},
				},
			},
			expectedError: "there is no skupper site in this namespace",
		},
		{
			name:  "link is not available",
			args:  []string{"my-link"},
			flags: common.CommandLinkUpdateFlags{Cost: "1", Timeout: time.Minute},
			skupperObjects: []runtime.Object{
				&v2alpha1.SiteList{
					Items: []v2alpha1.Site{
						{
							ObjectMeta: v1.ObjectMeta{
								Name:      "the-site",
								Namespace: "test",
							},
						},
					},
				},
			},
			expectedError: "the link \"my-link\" is not available in the namespace: links.skupper.io \"my-link\" not found",
		},
		{
			name:  "selected link does not exist",
			args:  []string{"my"},
			flags: common.CommandLinkUpdateFlags{Cost: "1", Timeout: time.Minute},
			skupperObjects: []runtime.Object{
				&v2alpha1.SiteList{
					Items: []v2alpha1.Site{
						{
							ObjectMeta: v1.ObjectMeta{
								Name:      "the-site",
								Namespace: "test",
							},
						},
					},
				},
				&v2alpha1.Link{
					ObjectMeta: v1.ObjectMeta{
						Name:      "my-link",
						Namespace: "test",
					},
				},
			},
			expectedError: "the link \"my\" is not available in the namespace: links.skupper.io \"my\" not found",
		},
		{
			name:  "link name is not specified.",
			args:  []string{},
			flags: common.CommandLinkUpdateFlags{Cost: "1", Timeout: time.Minute},
			skupperObjects: []runtime.Object{
				&v2alpha1.SiteList{
					Items: []v2alpha1.Site{
						{
							ObjectMeta: v1.ObjectMeta{
								Name:      "the-site",
								Namespace: "test",
							},
						},
					},
				},
				&v2alpha1.Link{
					ObjectMeta: v1.ObjectMeta{
						Name:      "my-link",
						Namespace: "test",
					},
				},
			},
			expectedError: "link name must not be empty",
		},
		{
			name:  "more than one argument was specified",
			args:  []string{"my", "link"},
			flags: common.CommandLinkUpdateFlags{Cost: "1", Timeout: time.Minute},
			skupperObjects: []runtime.Object{
				&v2alpha1.SiteList{
					Items: []v2alpha1.Site{
						{
							ObjectMeta: v1.ObjectMeta{
								Name:      "the-site",
								Namespace: "test",
							},
						},
					},
				},
				&v2alpha1.Link{
					ObjectMeta: v1.ObjectMeta{
						Name:      "my-link",
						Namespace: "test",
					},
				},
			},
			expectedError: "only one argument is allowed for this command",
		},
		{
			name:  "Cost is not valid.",
			args:  []string{"my-link"},
			flags: common.CommandLinkUpdateFlags{Cost: "one", Timeout: time.Minute},
			skupperObjects: []runtime.Object{
				&v2alpha1.SiteList{
					Items: []v2alpha1.Site{
						{
							ObjectMeta: v1.ObjectMeta{
								Name:      "the-site",
								Namespace: "test",
							},
						},
					},
				},
				&v2alpha1.Link{
					ObjectMeta: v1.ObjectMeta{
						Name:      "my-link",
						Namespace: "test",
					},
				},
			},
			expectedError: "link cost is not valid: strconv.Atoi: parsing \"one\": invalid syntax",
		},
		{
			name:  "Cost is not positive",
			args:  []string{"my-link"},
			flags: common.CommandLinkUpdateFlags{Cost: "-4", Timeout: time.Minute},
			skupperObjects: []runtime.Object{
				&v2alpha1.SiteList{
					Items: []v2alpha1.Site{
						{
							ObjectMeta: v1.ObjectMeta{
								Name:      "the-site",
								Namespace: "test",
							},
						},
					},
				},
				&v2alpha1.Link{
					ObjectMeta: v1.ObjectMeta{
						Name:      "my-link",
						Namespace: "test",
					},
				},
			},
			expectedError: "link cost is not valid: value is not positive",
		},
		{
			name:  "tls secret not available",
			args:  []string{"my-link"},
			flags: common.CommandLinkUpdateFlags{Cost: "1", TlsCredentials: "secret", Timeout: time.Minute},
			skupperObjects: []runtime.Object{
				&v2alpha1.SiteList{
					Items: []v2alpha1.Site{
						{
							ObjectMeta: v1.ObjectMeta{
								Name:      "the-site",
								Namespace: "test",
							},
						},
					},
				},
				&v2alpha1.Link{
					ObjectMeta: v1.ObjectMeta{
						Name:      "my-link",
						Namespace: "test",
					},
				},
			},
			expectedError: "the TLS secret \"secret\" is not available in the namespace: secrets \"secret\" not found",
		},
		{
			name:  "Timeout value is 0",
			args:  []string{"my-link"},
			flags: common.CommandLinkUpdateFlags{Cost: "1", TlsCredentials: "secret", Timeout: time.Second * 0},
			k8sObjects: []runtime.Object{
				&v12.Secret{
					ObjectMeta: v1.ObjectMeta{
						Name:      "secret",
						Namespace: "test",
					},
				},
			},
			skupperObjects: []runtime.Object{
				&v2alpha1.SiteList{
					Items: []v2alpha1.Site{
						{
							ObjectMeta: v1.ObjectMeta{
								Name:      "the-site",
								Namespace: "test",
							},
						},
					},
				},
				&v2alpha1.Link{
					ObjectMeta: v1.ObjectMeta{
						Name:      "my-link",
						Namespace: "test",
					},
				},
			},
			expectedError: "timeout is not valid: duration must not be less than 10s; got 0s",
		},
		{
			name:       "wait status is not valid",
			args:       []string{"east-link"},
			flags:      common.CommandLinkUpdateFlags{Cost: "1", Timeout: time.Minute, Wait: "created"},
			k8sObjects: nil,
			skupperObjects: []runtime.Object{
				&v2alpha1.SiteList{
					Items: []v2alpha1.Site{
						{
							ObjectMeta: v1.ObjectMeta{
								Name:      "the-site",
								Namespace: "test",
							},
						},
					},
				},
				&v2alpha1.Link{
					ObjectMeta: v1.ObjectMeta{
						Name:      "east-link",
						Namespace: "test",
					},
				},
			},
			expectedError: "status is not valid: value created not allowed. It should be one of this options: [ready configured none]",
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {

			command, err := newCmdLinkUpdateWithMocks("test", test.k8sObjects, test.skupperObjects, test.skupperError)
			assert.Assert(t, err)

			command.Flags = &test.flags

			testutils.CheckValidateInput(t, command, test.expectedError, test.args)
		})
	}
}

func TestCmdLinkUpdate_InputToOptions(t *testing.T) {

	type test struct {
		name                   string
		args                   []string
		flags                  common.CommandLinkUpdateFlags
		expectedTlsCredentials string
		expectedCost           int
		expectedTimeout        time.Duration
		expectedStatus         string
	}

	testTable := []test{
		{
			name:                   "check options",
			args:                   []string{"my-link"},
			flags:                  common.CommandLinkUpdateFlags{TlsCredentials: "secret", Cost: "1", Timeout: time.Minute, Wait: "ready"},
			expectedCost:           1,
			expectedTlsCredentials: "secret",
			expectedTimeout:        time.Minute,
			expectedStatus:         "ready",
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {

			cmd, err := newCmdLinkUpdateWithMocks("test", nil, nil, "")
			assert.Assert(t, err)
			cmd.Flags = &test.flags

			cmd.InputToOptions()

			assert.Check(t, cmd.tlsCredentials == test.expectedTlsCredentials)
			assert.Check(t, cmd.cost == test.expectedCost)
			assert.Check(t, cmd.timeout == test.expectedTimeout)
			assert.Check(t, cmd.status == test.expectedStatus)
		})
	}
}

func TestCmdLinkUpdate_Run(t *testing.T) {
	type test struct {
		name                string
		k8sObjects          []runtime.Object
		skupperObjects      []runtime.Object
		linkName            string
		Cost                int
		tlsCredentials      string
		errorMessage        string
		skupperErrorMessage string
	}

	testTable := []test{
		{
			name:           "runs ok",
			linkName:       "my-link",
			Cost:           1,
			tlsCredentials: "secret",
			skupperObjects: []runtime.Object{
				&v2alpha1.Link{
					ObjectMeta: v1.ObjectMeta{
						Name:      "my-link",
						Namespace: "test",
					},
				},
			},
		},
		{
			name:                "run fails",
			linkName:            "my-link",
			skupperErrorMessage: "error",
			errorMessage:        "error",
		},
		{
			name:         "run fails because link does not exist",
			linkName:     "my-link",
			errorMessage: "links.skupper.io \"my-link\" not found",
		},
		{
			name:     "runs ok without updating link",
			linkName: "my-link",
			skupperObjects: []runtime.Object{
				&v2alpha1.Link{
					ObjectMeta: v1.ObjectMeta{
						Name:      "my-link",
						Namespace: "test",
					},
				},
			},
		},
	}

	for _, test := range testTable {
		cmd, err := newCmdLinkUpdateWithMocks("test", test.k8sObjects, test.skupperObjects, test.skupperErrorMessage)
		assert.Assert(t, err)

		cmd.linkName = test.linkName
		cmd.tlsCredentials = test.tlsCredentials
		cmd.cost = test.Cost

		t.Run(test.name, func(t *testing.T) {

			err := cmd.Run()
			if err != nil {
				assert.Check(t, test.errorMessage == err.Error(), err.Error())
			} else {
				assert.Check(t, err == nil)
			}
		})
	}
}

func TestCmdLinkUpdate_WaitUntil(t *testing.T) {
	type test struct {
		name                string
		status              string
		k8sObjects          []runtime.Object
		skupperObjects      []runtime.Object
		skupperErrorMessage string
		linkName            string
		timeout             time.Duration
		expectError         bool
	}

	testTable := []test{
		{
			name:   "link is not configured",
			status: "configured",
			skupperObjects: []runtime.Object{
				&v2alpha1.Link{
					ObjectMeta: v1.ObjectMeta{
						Name:      "my-link",
						Namespace: "test",
					},
					Status: v2alpha1.LinkStatus{},
				},
			},
			linkName:    "my-link",
			timeout:     time.Second,
			expectError: true,
		},
		{
			name:        "link is not returned",
			linkName:    "my-link",
			timeout:     time.Second,
			expectError: true,
		},
		{
			name:   "link is ready",
			status: "ready",
			skupperObjects: []runtime.Object{
				&v2alpha1.Link{
					ObjectMeta: v1.ObjectMeta{
						Name:      "my-link",
						Namespace: "test",
					},
					Status: v2alpha1.LinkStatus{
						Status: v2alpha1.Status{
							Message: "OK",
							Conditions: []v1.Condition{
								{
									Message:            "OK",
									ObservedGeneration: 1,
									Reason:             "OK",
									Status:             "True",
									Type:               "Ready",
								},
							},
						},
					},
				},
			},
			linkName:    "my-link",
			timeout:     time.Second,
			expectError: false,
		},
		{
			name:       "link is not ready yet, but user waits for configured",
			status:     "configured",
			timeout:    time.Second,
			k8sObjects: nil,
			skupperObjects: []runtime.Object{
				&v2alpha1.Link{
					ObjectMeta: v1.ObjectMeta{
						Name:      "my-link",
						Namespace: "test",
					},
					Status: v2alpha1.LinkStatus{
						Status: v2alpha1.Status{
							Message: "OK",
							Conditions: []v1.Condition{
								{
									Message:            "OK",
									ObservedGeneration: 1,
									Reason:             "OK",
									Status:             "True",
									Type:               "Configured",
								},
							},
						},
					},
				},
			},
			linkName:    "my-link",
			expectError: false,
		},
		{
			name:       "user does not wait",
			status:     "none",
			timeout:    time.Second,
			k8sObjects: nil,
			skupperObjects: []runtime.Object{
				&v2alpha1.Link{
					ObjectMeta: v1.ObjectMeta{
						Name:      "my-link",
						Namespace: "test",
					},
					Status: v2alpha1.LinkStatus{
						Status: v2alpha1.Status{
							Message: "OK",
							Conditions: []v1.Condition{
								{
									Message:            "OK",
									ObservedGeneration: 1,
									Reason:             "OK",
									Status:             "True",
									Type:               "Configured",
								},
							},
						},
					},
				},
			},
			linkName:    "my-link",
			expectError: false,
		},
		{
			name:       "user waits for configured, but link had some errors while being configured",
			status:     "configured",
			timeout:    time.Second,
			k8sObjects: nil,
			skupperObjects: []runtime.Object{
				&v2alpha1.Link{
					ObjectMeta: v1.ObjectMeta{
						Name:      "my-link",
						Namespace: "test",
					},
					Status: v2alpha1.LinkStatus{
						Status: v2alpha1.Status{
							Message: "OK",
							Conditions: []v1.Condition{
								{
									Message:            "Error",
									ObservedGeneration: 1,
									Reason:             "Error",
									Status:             "False",
									Type:               "Configured",
								},
							},
						},
					},
				},
			},
			linkName:    "my-link",
			expectError: true,
		},
	}

	for _, test := range testTable {
		cmd, err := newCmdLinkUpdateWithMocks("test", test.k8sObjects, test.skupperObjects, test.skupperErrorMessage)
		assert.Assert(t, err)
		cmd.linkName = test.linkName
		cmd.timeout = test.timeout
		cmd.status = test.status

		t.Run(test.name, func(t *testing.T) {

			err := cmd.WaitUntil()

			if test.expectError {
				assert.Check(t, err != nil)
			} else {
				assert.Assert(t, err)
			}

		})
	}
}

// --- helper methods

func newCmdLinkUpdateWithMocks(namespace string, k8sObjects []runtime.Object, skupperObjects []runtime.Object, fakeSkupperError string) (*CmdLinkUpdate, error) {

	// We make sure the interval is appropriate
	utils.SetRetryProfile(utils.TestRetryProfile)

	client, err := fakeclient.NewFakeClient(namespace, k8sObjects, skupperObjects, fakeSkupperError)
	if err != nil {
		return nil, err
	}
	cmdLinkUpdate := &CmdLinkUpdate{
		Client:     client.GetSkupperClient().SkupperV2alpha1(),
		KubeClient: client.GetKubeClient(),
		Namespace:  namespace,
	}

	return cmdLinkUpdate, nil
}
