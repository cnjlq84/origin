package v1

import (
	"k8s.io/kubernetes/pkg/api/unversioned"
	kapi "k8s.io/kubernetes/pkg/api/v1"
	kruntime "k8s.io/kubernetes/pkg/runtime"
)

// Authorization is calculated against
// 1. all deny RoleBinding PolicyRules in the master namespace - short circuit on match
// 2. all allow RoleBinding PolicyRules in the master namespace - short circuit on match
// 3. all deny RoleBinding PolicyRules in the namespace - short circuit on match
// 4. all allow RoleBinding PolicyRules in the namespace - short circuit on match
// 5. deny by default

// PolicyRule holds information that describes a policy rule, but does not contain information
// about who the rule applies to or which namespace the rule applies to.
type PolicyRule struct {
	// Verbs is a list of Verbs that apply to ALL the ResourceKinds and AttributeRestrictions contained in this rule.  VerbAll represents all kinds.
	Verbs []string `json:"verbs"`
	// AttributeRestrictions will vary depending on what the Authorizer/AuthorizationAttributeBuilder pair supports.
	// If the Authorizer does not recognize how to handle the AttributeRestrictions, the Authorizer should report an error.
	AttributeRestrictions kruntime.RawExtension `json:"attributeRestrictions,omitempty"`
	// APIGroups is the name of the APIGroup that contains the resources.  If this field is empty, then both kubernetes and origin API groups are assumed.
	// That means that if an action is requested against one of the enumerated resources in either the kubernetes or the origin API group, the request
	// will be allowed
	APIGroups []string `json:"apiGroups"`
	// Resources is a list of resources this rule applies to.  ResourceAll represents all resources.
	Resources []string `json:"resources"`
	// ResourceNames is an optional white list of names that the rule applies to.  An empty set means that everything is allowed.
	ResourceNames []string `json:"resourceNames,omitempty"`
	// NonResourceURLsSlice is a set of partial urls that a user should have access to.  *s are allowed, but only as the full, final step in the path
	// This name is intentionally different than the internal type so that the DefaultConvert works nicely and because the ordering may be different.
	NonResourceURLsSlice []string `json:"nonResourceURLs,omitempty"`
}

// IsPersonalSubjectAccessReview is a marker for PolicyRule.AttributeRestrictions that denotes that subjectaccessreviews on self should be allowed
type IsPersonalSubjectAccessReview struct {
	unversioned.TypeMeta `json:",inline"`
}

// Role is a logical grouping of PolicyRules that can be referenced as a unit by RoleBindings.
type Role struct {
	unversioned.TypeMeta `json:",inline"`
	// Standard object's metadata.
	kapi.ObjectMeta `json:"metadata,omitempty"`

	// Rules holds all the PolicyRules for this Role
	Rules []PolicyRule `json:"rules"`
}

// RoleBinding references a Role, but not contain it.  It can reference any Role in the same namespace or in the global namespace.
// It adds who information via Users and Groups and namespace information by which namespace it exists in.  RoleBindings in a given
// namespace only have effect in that namespace (excepting the master namespace which has power in all namespaces).
type RoleBinding struct {
	unversioned.TypeMeta `json:",inline"`
	// Standard object's metadata.
	kapi.ObjectMeta `json:"metadata,omitempty"`

	// UserNames holds all the usernames directly bound to the role
	UserNames []string `json:"userNames"`
	// GroupNames holds all the groups directly bound to the role
	GroupNames []string `json:"groupNames"`
	// Subjects hold object references to authorize with this rule
	Subjects []kapi.ObjectReference `json:"subjects"`

	// RoleRef can only reference the current namespace and the global namespace
	// If the RoleRef cannot be resolved, the Authorizer must return an error.
	// Since Policy is a singleton, this is sufficient knowledge to locate a role
	RoleRef kapi.ObjectReference `json:"roleRef"`
}

// +genclient=true

// Policy is a object that holds all the Roles for a particular namespace.  There is at most
// one Policy document per namespace.
type Policy struct {
	unversioned.TypeMeta `json:",inline"`
	// Standard object's metadata.
	kapi.ObjectMeta `json:"metadata,omitempty"`

	// LastModified is the last time that any part of the Policy was created, updated, or deleted
	LastModified unversioned.Time `json:"lastModified"`

	// Roles holds all the Roles held by this Policy, mapped by Role.Name
	Roles []NamedRole `json:"roles"`
}

// PolicyBinding is a object that holds all the RoleBindings for a particular namespace.  There is
// one PolicyBinding document per referenced Policy namespace
type PolicyBinding struct {
	unversioned.TypeMeta `json:",inline"`
	// Standard object's metadata.
	kapi.ObjectMeta `json:"metadata,omitempty"`

	// LastModified is the last time that any part of the PolicyBinding was created, updated, or deleted
	LastModified unversioned.Time `json:"lastModified"`

	// PolicyRef is a reference to the Policy that contains all the Roles that this PolicyBinding's RoleBindings may reference
	PolicyRef kapi.ObjectReference `json:"policyRef"`
	// RoleBindings holds all the RoleBindings held by this PolicyBinding, mapped by RoleBinding.Name
	RoleBindings []NamedRoleBinding `json:"roleBindings"`
}

// NamedRole relates a Role with a name
type NamedRole struct {
	// Name is the name of the role
	Name string `json:"name"`
	// Role is the role being named
	Role Role `json:"role"`
}

// NamedRoleBinding relates a role binding with a name
type NamedRoleBinding struct {
	// Name is the name of the role binding
	Name string `json:"name"`
	// RoleBinding is the role binding being named
	RoleBinding RoleBinding `json:"roleBinding"`
}

// SelfSubjectRulesReview is a resource you can create to determine which actions you can perform in a namespace
type SelfSubjectRulesReview struct {
	unversioned.TypeMeta `json:",inline"`

	// Spec adds information about how to conduct the check
	Spec SelfSubjectRulesReviewSpec `json:"spec"`

	// Status is completed by the server to tell which permissions you have
	Status SubjectRulesReviewStatus `json:"status,omitempty"`
}

// SelfSubjectRulesReviewSpec adds information about how to conduct the check
type SelfSubjectRulesReviewSpec struct {
	// Scopes to use for the evaluation.  Empty means "use the unscoped (full) permissions of the user/groups".
	// Nil means "use the scopes on this request".
	Scopes []string `json:"scopes"`
}

// SubjectRulesReviewStatus is contains the result of a rules check
type SubjectRulesReviewStatus struct {
	// Rules is the list of rules (no particular sort) that are allowed for the subject
	Rules []PolicyRule `json:"rules"`
	// EvaluationError can appear in combination with Rules.  It means some error happened during evaluation
	// that may have prevented additional rules from being populated.
	EvaluationError string `json:"evaluationError,omitempty"`
}

// ResourceAccessReviewResponse describes who can perform the action
type ResourceAccessReviewResponse struct {
	unversioned.TypeMeta `json:",inline"`

	// Namespace is the namespace used for the access review
	Namespace string `json:"namespace,omitempty"`
	// UsersSlice is the list of users who can perform the action
	UsersSlice []string `json:"users"`
	// GroupsSlice is the list of groups who can perform the action
	GroupsSlice []string `json:"groups"`
}

// ResourceAccessReview is a means to request a list of which users and groups are authorized to perform the
// action specified by spec
type ResourceAccessReview struct {
	unversioned.TypeMeta `json:",inline"`

	// AuthorizationAttributes describes the action being tested.
	AuthorizationAttributes `json:",inline"`
}

// SubjectAccessReviewResponse describes whether or not a user or group can perform an action
type SubjectAccessReviewResponse struct {
	unversioned.TypeMeta `json:",inline"`

	// Namespace is the namespace used for the access review
	Namespace string `json:"namespace,omitempty"`
	// Allowed is required.  True if the action would be allowed, false otherwise.
	Allowed bool `json:"allowed"`
	// Reason is optional.  It indicates why a request was allowed or denied.
	Reason string `json:"reason,omitempty"`
}

// SubjectAccessReview is an object for requesting information about whether a user or group can perform an action
type SubjectAccessReview struct {
	unversioned.TypeMeta `json:",inline"`

	// AuthorizationAttributes describes the action being tested.
	AuthorizationAttributes `json:",inline"`
	// User is optional. If both User and Groups are empty, the current authenticated user is used.
	User string `json:"user"`
	// GroupsSlice is optional. Groups is the list of groups to which the User belongs.
	GroupsSlice []string `json:"groups"`
	// Scopes to use for the evaluation.  Empty means "use the unscoped (full) permissions of the user/groups".
	// Nil for a self-SAR, means "use the scopes on this request".
	// Nil for a regular SAR, means the same as empty.
	Scopes []string `json:"scopes"`
}

// LocalResourceAccessReview is a means to request a list of which users and groups are authorized to perform the action specified by spec in a particular namespace
type LocalResourceAccessReview struct {
	unversioned.TypeMeta `json:",inline"`

	// AuthorizationAttributes describes the action being tested.  The Namespace element is FORCED to the current namespace.
	AuthorizationAttributes `json:",inline"`
}

// LocalSubjectAccessReview is an object for requesting information about whether a user or group can perform an action in a particular namespace
type LocalSubjectAccessReview struct {
	unversioned.TypeMeta `json:",inline"`

	// AuthorizationAttributes describes the action being tested.  The Namespace element is FORCED to the current namespace.
	AuthorizationAttributes `json:",inline"`
	// User is optional.  If both User and Groups are empty, the current authenticated user is used.
	User string `json:"user"`
	// Groups is optional.  Groups is the list of groups to which the User belongs.
	GroupsSlice []string `json:"groups"`
	// Scopes to use for the evaluation.  Empty means "use the unscoped (full) permissions of the user/groups".
	// Nil for a self-SAR, means "use the scopes on this request".
	// Nil for a regular SAR, means the same as empty.
	Scopes []string `json:"scopes"`
}

// AuthorizationAttributes describes a request to the API server
type AuthorizationAttributes struct {
	// Namespace is the namespace of the action being requested.  Currently, there is no distinction between no namespace and all namespaces
	Namespace string `json:"namespace"`
	// Verb is one of: get, list, watch, create, update, delete
	Verb string `json:"verb"`
	// Group is the API group of the resource
	// Serialized as resourceAPIGroup to avoid confusion with the 'groups' field when inlined
	Group string `json:"resourceAPIGroup"`
	// Version is the API version of the resource
	// Serialized as resourceAPIVersion to avoid confusion with TypeMeta.apiVersion and ObjectMeta.resourceVersion when inlined
	Version string `json:"resourceAPIVersion"`
	// Resource is one of the existing resource types
	Resource string `json:"resource"`
	// ResourceName is the name of the resource being requested for a "get" or deleted for a "delete"
	ResourceName string `json:"resourceName"`
	// Content is the actual content of the request for create and update
	Content kruntime.RawExtension `json:"content,omitempty"`
}

// PolicyList is a collection of Policies
type PolicyList struct {
	unversioned.TypeMeta `json:",inline"`
	// Standard object's metadata.
	unversioned.ListMeta `json:"metadata,omitempty"`

	// Items is a list of Policies
	Items []Policy `json:"items"`
}

// PolicyBindingList is a collection of PolicyBindings
type PolicyBindingList struct {
	unversioned.TypeMeta `json:",inline"`
	// Standard object's metadata.
	unversioned.ListMeta `json:"metadata,omitempty"`

	// Items is a list of PolicyBindings
	Items []PolicyBinding `json:"items"`
}

// RoleBindingList is a collection of RoleBindings
type RoleBindingList struct {
	unversioned.TypeMeta `json:",inline"`
	// Standard object's metadata.
	unversioned.ListMeta `json:"metadata,omitempty"`

	// Items is a list of RoleBindings
	Items []RoleBinding `json:"items"`
}

// RoleList is a collection of Roles
type RoleList struct {
	unversioned.TypeMeta `json:",inline"`
	// Standard object's metadata.
	unversioned.ListMeta `json:"metadata,omitempty"`

	// Items is a list of Roles
	Items []Role `json:"items"`
}

// ClusterRole is a logical grouping of PolicyRules that can be referenced as a unit by ClusterRoleBindings.
type ClusterRole struct {
	unversioned.TypeMeta `json:",inline"`
	// Standard object's metadata.
	kapi.ObjectMeta `json:"metadata,omitempty"`

	// Rules holds all the PolicyRules for this ClusterRole
	Rules []PolicyRule `json:"rules"`
}

// ClusterRoleBinding references a ClusterRole, but not contain it.  It can reference any ClusterRole in the same namespace or in the global namespace.
// It adds who information via Users and Groups and namespace information by which namespace it exists in.  ClusterRoleBindings in a given
// namespace only have effect in that namespace (excepting the master namespace which has power in all namespaces).
type ClusterRoleBinding struct {
	unversioned.TypeMeta `json:",inline"`
	// Standard object's metadata.
	kapi.ObjectMeta `json:"metadata,omitempty"`

	// UserNames holds all the usernames directly bound to the role
	UserNames []string `json:"userNames"`
	// GroupNames holds all the groups directly bound to the role
	GroupNames []string `json:"groupNames"`
	// Subjects hold object references to authorize with this rule
	Subjects []kapi.ObjectReference `json:"subjects"`

	// RoleRef can only reference the current namespace and the global namespace
	// If the ClusterRoleRef cannot be resolved, the Authorizer must return an error.
	// Since Policy is a singleton, this is sufficient knowledge to locate a role
	RoleRef kapi.ObjectReference `json:"roleRef"`
}

// ClusterPolicy is a object that holds all the ClusterRoles for a particular namespace.  There is at most
// one ClusterPolicy document per namespace.
type ClusterPolicy struct {
	unversioned.TypeMeta `json:",inline"`
	// Standard object's metadata.
	kapi.ObjectMeta `json:"metadata,omitempty"`

	// LastModified is the last time that any part of the ClusterPolicy was created, updated, or deleted
	LastModified unversioned.Time `json:"lastModified"`

	// Roles holds all the ClusterRoles held by this ClusterPolicy, mapped by ClusterRole.Name
	Roles []NamedClusterRole `json:"roles"`
}

// ClusterPolicyBinding is a object that holds all the ClusterRoleBindings for a particular namespace.  There is
// one ClusterPolicyBinding document per referenced ClusterPolicy namespace
type ClusterPolicyBinding struct {
	unversioned.TypeMeta `json:",inline"`
	// Standard object's metadata.
	kapi.ObjectMeta `json:"metadata,omitempty"`

	// LastModified is the last time that any part of the ClusterPolicyBinding was created, updated, or deleted
	LastModified unversioned.Time `json:"lastModified"`

	// PolicyRef is a reference to the ClusterPolicy that contains all the ClusterRoles that this ClusterPolicyBinding's RoleBindings may reference
	PolicyRef kapi.ObjectReference `json:"policyRef"`
	// RoleBindings holds all the ClusterRoleBindings held by this ClusterPolicyBinding, mapped by ClusterRoleBinding.Name
	RoleBindings []NamedClusterRoleBinding `json:"roleBindings"`
}

// NamedClusterRole relates a name with a cluster role
type NamedClusterRole struct {
	// Name is the name of the cluster role
	Name string `json:"name"`
	// Role is the cluster role being named
	Role ClusterRole `json:"role"`
}

// NamedClusterRoleBinding relates a name with a cluster role binding
type NamedClusterRoleBinding struct {
	// Name is the name of the cluster role binding
	Name string `json:"name"`
	// RoleBinding is the cluster role binding being named
	RoleBinding ClusterRoleBinding `json:"roleBinding"`
}

// ClusterPolicyList is a collection of ClusterPolicies
type ClusterPolicyList struct {
	unversioned.TypeMeta `json:",inline"`
	// Standard object's metadata.
	unversioned.ListMeta `json:"metadata,omitempty"`

	// Items is a list of ClusterPolicies
	Items []ClusterPolicy `json:"items"`
}

// ClusterPolicyBindingList is a collection of ClusterPolicyBindings
type ClusterPolicyBindingList struct {
	unversioned.TypeMeta `json:",inline"`
	// Standard object's metadata.
	unversioned.ListMeta `json:"metadata,omitempty"`

	// Items is a list of ClusterPolicyBindings
	Items []ClusterPolicyBinding `json:"items"`
}

// ClusterRoleBindingList is a collection of ClusterRoleBindings
type ClusterRoleBindingList struct {
	unversioned.TypeMeta `json:",inline"`
	// Standard object's metadata.
	unversioned.ListMeta `json:"metadata,omitempty"`

	// Items is a list of ClusterRoleBindings
	Items []ClusterRoleBinding `json:"items"`
}

// ClusterRoleList is a collection of ClusterRoles
type ClusterRoleList struct {
	unversioned.TypeMeta `json:",inline"`
	// Standard object's metadata.
	unversioned.ListMeta `json:"metadata,omitempty"`

	// Items is a list of ClusterRoles
	Items []ClusterRole `json:"items"`
}
