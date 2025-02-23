// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/schematics-go-sdk/schematicsv1"
)

func dataSourceIBMSchematicsAction() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMSchematicsActionRead,

		Schema: map[string]*schema.Schema{
			"action_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Action Id.  Use GET /actions API to look up the Action Ids in your IBM Cloud account.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique name of your action. The name can be up to 128 characters long and can include alphanumeric characters, spaces, dashes, and underscores. **Example** you can use the name to stop action.",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Action description.",
			},
			"location": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "List of locations supported by IBM Cloud Schematics service.  While creating your workspace or action, choose the right region, since it cannot be changed.  Note, this does not limit the location of the IBM Cloud resources, provisioned using Schematics.",
			},
			"resource_group": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource-group name for an action.  By default, action is created in default resource group.",
			},
			"tags": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Action tags.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"user_state": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "User defined status of the Schematics object.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"state": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "User-defined states  * `draft` Object can be modified; can be used by Jobs run by the author, during execution  * `live` Object can be modified; can be used by Jobs during execution  * `locked` Object cannot be modified; can be used by Jobs during execution  * `disable` Object can be modified. cannot be used by Jobs during execution.",
						},
						"set_by": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the User who set the state of the Object.",
						},
						"set_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "When the User who set the state of the Object.",
						},
					},
				},
			},
			"source_readme_url": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL of the `README` file, for the source URL.",
			},
			"source": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Source of templates, playbooks, or controls.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of source for the Template.",
						},
						"git": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Connection details to Git source.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"computed_git_repo_url": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The Complete URL which is computed by git_repo_url, git_repo_folder and branch.",
									},
									"git_repo_url": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "URL to the GIT Repo that can be used to clone the template.",
									},
									"git_token": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Personal Access Token to connect to Git URLs.",
									},
									"git_repo_folder": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name of the folder in the Git Repo, that contains the template.",
									},
									"git_release": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name of the release tag, used to fetch the Git Repo.",
									},
									"git_branch": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name of the branch, used to fetch the Git Repo.",
									},
								},
							},
						},
						"catalog": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Connection details to IBM Cloud Catalog source.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"catalog_name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "name of the private catalog.",
									},
									"offering_name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name of the offering in the IBM Catalog.",
									},
									"offering_version": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Version string of the offering in the IBM Catalog.",
									},
									"offering_kind": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Type of the offering, in the IBM Catalog.",
									},
									"offering_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Id of the offering the IBM Catalog.",
									},
									"offering_version_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Id of the offering version the IBM Catalog.",
									},
									"offering_repo_url": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Repo Url of the offering, in the IBM Catalog.",
									},
								},
							},
						},
						"cos_bucket": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Connection details to a IBM Cloud Object Storage bucket.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cos_bucket_url": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "COS Bucket Url.",
									},
								},
							},
						},
					},
				},
			},
			"source_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of source for the Template.",
			},
			"command_parameter": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Schematics job command parameter (playbook-name).",
			},
			"inventory": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Target inventory record ID, used by the action or ansible playbook.",
			},
			"credentials": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "credentials of the Action.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the variable.",
						},
						"value": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Value for the variable or reference to the value.",
						},
						"metadata": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "User editable metadata for the variables.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Type of the variable.",
									},
									"aliases": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "List of aliases for the variable name.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"description": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Description of the meta data.",
									},
									"default_value": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Default value for the variable, if the override value is not specified.",
									},
									"secure": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Is the variable secure or sensitive ?.",
									},
									"immutable": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Is the variable readonly ?.",
									},
									"hidden": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "If true, the variable will not be displayed on UI or CLI.",
									},
									"options": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "List of possible values for this variable.  If type is integer or date, then the array of string will be  converted to array of integers or date during runtime.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"min_value": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Minimum value of the variable. Applicable for integer type.",
									},
									"max_value": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Maximum value of the variable. Applicable for integer type.",
									},
									"min_length": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Minimum length of the variable value. Applicable for string type.",
									},
									"max_length": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Maximum length of the variable value. Applicable for string type.",
									},
									"matches": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Regex for the variable value.",
									},
									"position": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Relative position of this variable in a list.",
									},
									"group_by": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Display name of the group this variable belongs to.",
									},
									"source": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Source of this meta-data.",
									},
								},
							},
						},
						"link": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Reference link to the variable value By default the expression will point to self.value.",
						},
					},
				},
			},
			"bastion": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Describes a bastion resource.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Bastion Name(Unique).",
						},
						"host": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Reference to the Inventory resource definition.",
						},
					},
				},
			},
			"bastion_credential": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "User editable variable data & system generated reference to value.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the variable.",
						},
						"value": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Value for the variable or reference to the value.",
						},
						"metadata": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "User editable metadata for the variables.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Type of the variable.",
									},
									"aliases": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "List of aliases for the variable name.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"description": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Description of the meta data.",
									},
									"default_value": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Default value for the variable, if the override value is not specified.",
									},
									"secure": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Is the variable secure or sensitive ?.",
									},
									"immutable": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Is the variable readonly ?.",
									},
									"hidden": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "If true, the variable will not be displayed on UI or CLI.",
									},
									"options": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "List of possible values for this variable.  If type is integer or date, then the array of string will be  converted to array of integers or date during runtime.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"min_value": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Minimum value of the variable. Applicable for integer type.",
									},
									"max_value": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Maximum value of the variable. Applicable for integer type.",
									},
									"min_length": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Minimum length of the variable value. Applicable for string type.",
									},
									"max_length": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Maximum length of the variable value. Applicable for string type.",
									},
									"matches": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Regex for the variable value.",
									},
									"position": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Relative position of this variable in a list.",
									},
									"group_by": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Display name of the group this variable belongs to.",
									},
									"source": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Source of this meta-data.",
									},
								},
							},
						},
						"link": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Reference link to the variable value By default the expression will point to self.value.",
						},
					},
				},
			},
			"targets_ini": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Inventory of host and host group for the playbook in `INI` file format. For example, `\"targets_ini\": \"[webserverhost]  172.22.192.6  [dbhost]  172.22.192.5\"`. For more information, about an inventory host group syntax, see [Inventory host groups](https://cloud.ibm.com/docs/schematics?topic=schematics-schematics-cli-reference#schematics-inventory-host-grps).",
			},
			"action_inputs": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Input variables for the Action.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the variable.",
						},
						"value": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Value for the variable or reference to the value.",
						},
						"metadata": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "User editable metadata for the variables.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Type of the variable.",
									},
									"aliases": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "List of aliases for the variable name.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"description": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Description of the meta data.",
									},
									"default_value": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Default value for the variable, if the override value is not specified.",
									},
									"secure": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Is the variable secure or sensitive ?.",
									},
									"immutable": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Is the variable readonly ?.",
									},
									"hidden": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "If true, the variable will not be displayed on UI or CLI.",
									},
									"options": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "List of possible values for this variable.  If type is integer or date, then the array of string will be  converted to array of integers or date during runtime.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"min_value": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Minimum value of the variable. Applicable for integer type.",
									},
									"max_value": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Maximum value of the variable. Applicable for integer type.",
									},
									"min_length": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Minimum length of the variable value. Applicable for string type.",
									},
									"max_length": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Maximum length of the variable value. Applicable for string type.",
									},
									"matches": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Regex for the variable value.",
									},
									"position": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Relative position of this variable in a list.",
									},
									"group_by": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Display name of the group this variable belongs to.",
									},
									"source": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Source of this meta-data.",
									},
								},
							},
						},
						"link": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Reference link to the variable value By default the expression will point to self.value.",
						},
					},
				},
			},
			"action_outputs": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Output variables for the Action.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the variable.",
						},
						"value": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Value for the variable or reference to the value.",
						},
						"metadata": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "User editable metadata for the variables.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Type of the variable.",
									},
									"aliases": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "List of aliases for the variable name.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"description": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Description of the meta data.",
									},
									"default_value": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Default value for the variable, if the override value is not specified.",
									},
									"secure": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Is the variable secure or sensitive ?.",
									},
									"immutable": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Is the variable readonly ?.",
									},
									"hidden": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "If true, the variable will not be displayed on UI or CLI.",
									},
									"options": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "List of possible values for this variable.  If type is integer or date, then the array of string will be  converted to array of integers or date during runtime.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"min_value": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Minimum value of the variable. Applicable for integer type.",
									},
									"max_value": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Maximum value of the variable. Applicable for integer type.",
									},
									"min_length": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Minimum length of the variable value. Applicable for string type.",
									},
									"max_length": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Maximum length of the variable value. Applicable for string type.",
									},
									"matches": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Regex for the variable value.",
									},
									"position": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Relative position of this variable in a list.",
									},
									"group_by": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Display name of the group this variable belongs to.",
									},
									"source": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Source of this meta-data.",
									},
								},
							},
						},
						"link": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Reference link to the variable value By default the expression will point to self.value.",
						},
					},
				},
			},
			"settings": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Environment variables for the Action.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the variable.",
						},
						"value": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Value for the variable or reference to the value.",
						},
						"metadata": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "User editable metadata for the variables.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Type of the variable.",
									},
									"aliases": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "List of aliases for the variable name.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"description": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Description of the meta data.",
									},
									"default_value": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Default value for the variable, if the override value is not specified.",
									},
									"secure": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Is the variable secure or sensitive ?.",
									},
									"immutable": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Is the variable readonly ?.",
									},
									"hidden": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "If true, the variable will not be displayed on UI or CLI.",
									},
									"options": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "List of possible values for this variable.  If type is integer or date, then the array of string will be  converted to array of integers or date during runtime.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"min_value": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Minimum value of the variable. Applicable for integer type.",
									},
									"max_value": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Maximum value of the variable. Applicable for integer type.",
									},
									"min_length": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Minimum length of the variable value. Applicable for string type.",
									},
									"max_length": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Maximum length of the variable value. Applicable for string type.",
									},
									"matches": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Regex for the variable value.",
									},
									"position": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Relative position of this variable in a list.",
									},
									"group_by": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Display name of the group this variable belongs to.",
									},
									"source": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Source of this meta-data.",
									},
								},
							},
						},
						"link": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Reference link to the variable value By default the expression will point to self.value.",
						},
					},
				},
			},
			"id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Action ID.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Action Cloud Resource Name.",
			},
			"account": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Action account ID.",
			},
			"source_created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Action Playbook Source creation time.",
			},
			"source_created_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "E-mail address of user who created the Action Playbook Source.",
			},
			"source_updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The action playbook updation time.",
			},
			"source_updated_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "E-mail address of user who updated the action playbook source.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Action creation time.",
			},
			"created_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "E-mail address of the user who created an action.",
			},
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Action updation time.",
			},
			"updated_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "E-mail address of the user who updated an action.",
			},
			"state": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Computed state of the Action.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status_code": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status of automation (workspace or action).",
						},
						"status_job_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Job id reference for this status.",
						},
						"status_message": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Automation status message - to be displayed along with the status_code.",
						},
					},
				},
			},
			"playbook_names": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Playbook names retrieved from the respository.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"sys_lock": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "System lock status.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sys_locked": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Is the automation locked by a Schematic job ?.",
						},
						"sys_locked_by": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the User who performed the job, that lead to the locking of the automation.",
						},
						"sys_locked_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "When the User performed the job that lead to locking of the automation ?.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMSchematicsActionRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	schematicsClient, err := meta.(ClientSession).SchematicsV1()
	if err != nil {
		return diag.FromErr(err)
	}

	getActionOptions := &schematicsv1.GetActionOptions{}

	getActionOptions.SetActionID(d.Get("action_id").(string))

	action, response, err := schematicsClient.GetActionWithContext(context, getActionOptions)
	if err != nil {
		log.Printf("[DEBUG] GetActionWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetActionWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s", *getActionOptions.ActionID))
	if err = d.Set("name", action.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if err = d.Set("description", action.Description); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting description: %s", err))
	}
	if err = d.Set("location", action.Location); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting location: %s", err))
	}
	if err = d.Set("resource_group", action.ResourceGroup); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_group: %s", err))
	}

	if action.UserState != nil {
		err = d.Set("user_state", dataSourceActionFlattenUserState(*action.UserState))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting user_state %s", err))
		}
	}
	if err = d.Set("source_readme_url", action.SourceReadmeURL); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting source_readme_url: %s", err))
	}

	if action.Source != nil {
		err = d.Set("source", dataSourceActionFlattenSource(*action.Source))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting source %s", err))
		}
	}
	if err = d.Set("source_type", action.SourceType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting source_type: %s", err))
	}
	if err = d.Set("command_parameter", action.CommandParameter); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting command_parameter: %s", err))
	}
	if err = d.Set("inventory", action.Inventory); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting inventory: %s", err))
	}

	if action.Credentials != nil {
		err = d.Set("credentials", dataSourceActionFlattenCredentials(action.Credentials))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting credentials %s", err))
		}
	}

	if action.Bastion != nil {
		err = d.Set("bastion", dataSourceActionFlattenBastion(*action.Bastion))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting bastion %s", err))
		}
	}

	if action.BastionCredential != nil {
		err = d.Set("bastion_credential", dataSourceActionFlattenBastionCredential(*action.BastionCredential))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting bastion_credential %s", err))
		}
	}
	if err = d.Set("targets_ini", action.TargetsIni); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting targets_ini: %s", err))
	}

	if action.Inputs != nil {
		err = d.Set("action_inputs", dataSourceActionFlattenInputs(action.Inputs))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting action_inputs %s", err))
		}
	}

	if action.Outputs != nil {
		err = d.Set("action_outputs", dataSourceActionFlattenOutputs(action.Outputs))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting action_outputs %s", err))
		}
	}

	if action.Settings != nil {
		err = d.Set("settings", dataSourceActionFlattenSettings(action.Settings))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting settings %s", err))
		}
	}
	if err = d.Set("id", action.ID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting id: %s", err))
	}
	if err = d.Set("crn", action.Crn); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting crn: %s", err))
	}
	if err = d.Set("account", action.Account); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting account: %s", err))
	}
	if err = d.Set("source_created_at", dateTimeToString(action.SourceCreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting source_created_at: %s", err))
	}
	if err = d.Set("source_created_by", action.SourceCreatedBy); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting source_created_by: %s", err))
	}
	if err = d.Set("source_updated_at", dateTimeToString(action.SourceUpdatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting source_updated_at: %s", err))
	}
	if err = d.Set("source_updated_by", action.SourceUpdatedBy); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting source_updated_by: %s", err))
	}
	if err = d.Set("created_at", dateTimeToString(action.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	if err = d.Set("created_by", action.CreatedBy); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_by: %s", err))
	}
	if err = d.Set("updated_at", dateTimeToString(action.UpdatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting updated_at: %s", err))
	}
	if err = d.Set("updated_by", action.UpdatedBy); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting updated_by: %s", err))
	}

	if action.State != nil {
		err = d.Set("state", dataSourceActionFlattenState(*action.State))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting state %s", err))
		}
	}

	if action.SysLock != nil {
		err = d.Set("sys_lock", dataSourceActionFlattenSysLock(*action.SysLock))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting sys_lock %s", err))
		}
	}

	return nil
}

func dataSourceActionFlattenUserState(result schematicsv1.UserState) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceActionUserStateToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceActionUserStateToMap(userStateItem schematicsv1.UserState) (userStateMap map[string]interface{}) {
	userStateMap = map[string]interface{}{}

	if userStateItem.State != nil {
		userStateMap["state"] = userStateItem.State
	}
	if userStateItem.SetBy != nil {
		userStateMap["set_by"] = userStateItem.SetBy
	}
	if userStateItem.SetAt != nil {
		userStateMap["set_at"] = userStateItem.SetAt.String()
	}

	return userStateMap
}

func dataSourceActionFlattenSource(result schematicsv1.ExternalSource) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceActionSourceToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceActionSourceToMap(sourceItem schematicsv1.ExternalSource) (sourceMap map[string]interface{}) {
	sourceMap = map[string]interface{}{}

	if sourceItem.SourceType != nil {
		sourceMap["source_type"] = sourceItem.SourceType
	}
	if sourceItem.Git != nil {
		gitList := []map[string]interface{}{}
		gitMap := dataSourceActionSourceGitToMap(*sourceItem.Git)
		gitList = append(gitList, gitMap)
		sourceMap["git"] = gitList
	}
	if sourceItem.Catalog != nil {
		catalogList := []map[string]interface{}{}
		catalogMap := dataSourceActionSourceCatalogToMap(*sourceItem.Catalog)
		catalogList = append(catalogList, catalogMap)
		sourceMap["catalog"] = catalogList
	}
	if sourceItem.CosBucket != nil {
		cosBucketList := []map[string]interface{}{}
		cosBucketMap := dataSourceActionSourceCosBucketToMap(*sourceItem.CosBucket)
		cosBucketList = append(cosBucketList, cosBucketMap)
		sourceMap["cos_bucket"] = cosBucketList
	}

	return sourceMap
}

func dataSourceActionSourceGitToMap(gitItem schematicsv1.ExternalSourceGit) (gitMap map[string]interface{}) {
	gitMap = map[string]interface{}{}

	if gitItem.ComputedGitRepoURL != nil {
		gitMap["computed_git_repo_url"] = gitItem.ComputedGitRepoURL
	}
	if gitItem.GitRepoURL != nil {
		gitMap["git_repo_url"] = gitItem.GitRepoURL
	}
	if gitItem.GitToken != nil {
		gitMap["git_token"] = gitItem.GitToken
	}
	if gitItem.GitRepoFolder != nil {
		gitMap["git_repo_folder"] = gitItem.GitRepoFolder
	}
	if gitItem.GitRelease != nil {
		gitMap["git_release"] = gitItem.GitRelease
	}
	if gitItem.GitBranch != nil {
		gitMap["git_branch"] = gitItem.GitBranch
	}

	return gitMap
}

func dataSourceActionSourceCatalogToMap(catalogItem schematicsv1.ExternalSourceCatalog) (catalogMap map[string]interface{}) {
	catalogMap = map[string]interface{}{}

	if catalogItem.CatalogName != nil {
		catalogMap["catalog_name"] = catalogItem.CatalogName
	}
	if catalogItem.OfferingName != nil {
		catalogMap["offering_name"] = catalogItem.OfferingName
	}
	if catalogItem.OfferingVersion != nil {
		catalogMap["offering_version"] = catalogItem.OfferingVersion
	}
	if catalogItem.OfferingKind != nil {
		catalogMap["offering_kind"] = catalogItem.OfferingKind
	}
	if catalogItem.OfferingID != nil {
		catalogMap["offering_id"] = catalogItem.OfferingID
	}
	if catalogItem.OfferingVersionID != nil {
		catalogMap["offering_version_id"] = catalogItem.OfferingVersionID
	}
	if catalogItem.OfferingRepoURL != nil {
		catalogMap["offering_repo_url"] = catalogItem.OfferingRepoURL
	}

	return catalogMap
}

func dataSourceActionSourceCosBucketToMap(cosBucketItem schematicsv1.ExternalSourceCosBucket) (cosBucketMap map[string]interface{}) {
	cosBucketMap = map[string]interface{}{}

	if cosBucketItem.CosBucketURL != nil {
		cosBucketMap["cos_bucket_url"] = cosBucketItem.CosBucketURL
	}

	return cosBucketMap
}

func dataSourceActionFlattenCredentials(result []schematicsv1.VariableData) (credentials []map[string]interface{}) {
	for _, credentialsItem := range result {
		credentials = append(credentials, dataSourceActionCredentialsToMap(credentialsItem))
	}

	return credentials
}

func dataSourceActionCredentialsToMap(credentialsItem schematicsv1.VariableData) (credentialsMap map[string]interface{}) {
	credentialsMap = map[string]interface{}{}

	if credentialsItem.Name != nil {
		credentialsMap["name"] = credentialsItem.Name
	}
	if credentialsItem.Value != nil {
		credentialsMap["value"] = credentialsItem.Value
	}
	if credentialsItem.Metadata != nil {
		metadataList := []map[string]interface{}{}
		metadataMap := dataSourceActionCredentialsMetadataToMap(*credentialsItem.Metadata)
		metadataList = append(metadataList, metadataMap)
		credentialsMap["metadata"] = metadataList
	}
	if credentialsItem.Link != nil {
		credentialsMap["link"] = credentialsItem.Link
	}

	return credentialsMap
}

func dataSourceActionCredentialsMetadataToMap(metadataItem schematicsv1.VariableMetadata) (metadataMap map[string]interface{}) {
	metadataMap = map[string]interface{}{}

	if metadataItem.Type != nil {
		metadataMap["type"] = metadataItem.Type
	}
	if metadataItem.Aliases != nil {
		metadataMap["aliases"] = metadataItem.Aliases
	}
	if metadataItem.Description != nil {
		metadataMap["description"] = metadataItem.Description
	}
	if metadataItem.DefaultValue != nil {
		metadataMap["default_value"] = metadataItem.DefaultValue
	}
	if metadataItem.Secure != nil {
		metadataMap["secure"] = metadataItem.Secure
	}
	if metadataItem.Immutable != nil {
		metadataMap["immutable"] = metadataItem.Immutable
	}
	if metadataItem.Hidden != nil {
		metadataMap["hidden"] = metadataItem.Hidden
	}
	if metadataItem.Options != nil {
		metadataMap["options"] = metadataItem.Options
	}
	if metadataItem.MinValue != nil {
		metadataMap["min_value"] = metadataItem.MinValue
	}
	if metadataItem.MaxValue != nil {
		metadataMap["max_value"] = metadataItem.MaxValue
	}
	if metadataItem.MinLength != nil {
		metadataMap["min_length"] = metadataItem.MinLength
	}
	if metadataItem.MaxLength != nil {
		metadataMap["max_length"] = metadataItem.MaxLength
	}
	if metadataItem.Matches != nil {
		metadataMap["matches"] = metadataItem.Matches
	}
	if metadataItem.Position != nil {
		metadataMap["position"] = metadataItem.Position
	}
	if metadataItem.GroupBy != nil {
		metadataMap["group_by"] = metadataItem.GroupBy
	}
	if metadataItem.Source != nil {
		metadataMap["source"] = metadataItem.Source
	}

	return metadataMap
}

func dataSourceActionFlattenBastion(result schematicsv1.BastionResourceDefinition) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceActionBastionToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceActionBastionToMap(bastionItem schematicsv1.BastionResourceDefinition) (bastionMap map[string]interface{}) {
	bastionMap = map[string]interface{}{}

	if bastionItem.Name != nil {
		bastionMap["name"] = bastionItem.Name
	}
	if bastionItem.Host != nil {
		bastionMap["host"] = bastionItem.Host
	}

	return bastionMap
}

func dataSourceActionFlattenBastionCredential(result schematicsv1.VariableData) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceActionBastionCredentialToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceActionBastionCredentialToMap(bastionCredentialItem schematicsv1.VariableData) (bastionCredentialMap map[string]interface{}) {
	bastionCredentialMap = map[string]interface{}{}

	if bastionCredentialItem.Name != nil {
		bastionCredentialMap["name"] = bastionCredentialItem.Name
	}
	if bastionCredentialItem.Value != nil {
		bastionCredentialMap["value"] = bastionCredentialItem.Value
	}
	if bastionCredentialItem.Metadata != nil {
		metadataList := []map[string]interface{}{}
		metadataMap := dataSourceActionBastionCredentialMetadataToMap(*bastionCredentialItem.Metadata)
		metadataList = append(metadataList, metadataMap)
		bastionCredentialMap["metadata"] = metadataList
	}
	if bastionCredentialItem.Link != nil {
		bastionCredentialMap["link"] = bastionCredentialItem.Link
	}

	return bastionCredentialMap
}

func dataSourceActionBastionCredentialMetadataToMap(metadataItem schematicsv1.VariableMetadata) (metadataMap map[string]interface{}) {
	metadataMap = map[string]interface{}{}

	if metadataItem.Type != nil {
		metadataMap["type"] = metadataItem.Type
	}
	if metadataItem.Aliases != nil {
		metadataMap["aliases"] = metadataItem.Aliases
	}
	if metadataItem.Description != nil {
		metadataMap["description"] = metadataItem.Description
	}
	if metadataItem.DefaultValue != nil {
		metadataMap["default_value"] = metadataItem.DefaultValue
	}
	if metadataItem.Secure != nil {
		metadataMap["secure"] = metadataItem.Secure
	}
	if metadataItem.Immutable != nil {
		metadataMap["immutable"] = metadataItem.Immutable
	}
	if metadataItem.Hidden != nil {
		metadataMap["hidden"] = metadataItem.Hidden
	}
	if metadataItem.Options != nil {
		metadataMap["options"] = metadataItem.Options
	}
	if metadataItem.MinValue != nil {
		metadataMap["min_value"] = metadataItem.MinValue
	}
	if metadataItem.MaxValue != nil {
		metadataMap["max_value"] = metadataItem.MaxValue
	}
	if metadataItem.MinLength != nil {
		metadataMap["min_length"] = metadataItem.MinLength
	}
	if metadataItem.MaxLength != nil {
		metadataMap["max_length"] = metadataItem.MaxLength
	}
	if metadataItem.Matches != nil {
		metadataMap["matches"] = metadataItem.Matches
	}
	if metadataItem.Position != nil {
		metadataMap["position"] = metadataItem.Position
	}
	if metadataItem.GroupBy != nil {
		metadataMap["group_by"] = metadataItem.GroupBy
	}
	if metadataItem.Source != nil {
		metadataMap["source"] = metadataItem.Source
	}

	return metadataMap
}

func dataSourceActionFlattenInputs(result []schematicsv1.VariableData) (actionInputs []map[string]interface{}) {
	for _, actionInputsItem := range result {
		actionInputs = append(actionInputs, dataSourceActionInputsToMap(actionInputsItem))
	}

	return actionInputs
}

func dataSourceActionInputsToMap(inputsItem schematicsv1.VariableData) (inputsMap map[string]interface{}) {
	inputsMap = map[string]interface{}{}

	if inputsItem.Name != nil {
		inputsMap["name"] = inputsItem.Name
	}
	if inputsItem.Value != nil {
		inputsMap["value"] = inputsItem.Value
	}
	if inputsItem.Metadata != nil {
		metadataList := []map[string]interface{}{}
		metadataMap := dataSourceActionInputsMetadataToMap(*inputsItem.Metadata)
		metadataList = append(metadataList, metadataMap)
		inputsMap["metadata"] = metadataList
	}
	if inputsItem.Link != nil {
		inputsMap["link"] = inputsItem.Link
	}

	return inputsMap
}

func dataSourceActionInputsMetadataToMap(metadataItem schematicsv1.VariableMetadata) (metadataMap map[string]interface{}) {
	metadataMap = map[string]interface{}{}

	if metadataItem.Type != nil {
		metadataMap["type"] = metadataItem.Type
	}
	if metadataItem.Aliases != nil {
		metadataMap["aliases"] = metadataItem.Aliases
	}
	if metadataItem.Description != nil {
		metadataMap["description"] = metadataItem.Description
	}
	if metadataItem.DefaultValue != nil {
		metadataMap["default_value"] = metadataItem.DefaultValue
	}
	if metadataItem.Secure != nil {
		metadataMap["secure"] = metadataItem.Secure
	}
	if metadataItem.Immutable != nil {
		metadataMap["immutable"] = metadataItem.Immutable
	}
	if metadataItem.Hidden != nil {
		metadataMap["hidden"] = metadataItem.Hidden
	}
	if metadataItem.Options != nil {
		metadataMap["options"] = metadataItem.Options
	}
	if metadataItem.MinValue != nil {
		metadataMap["min_value"] = metadataItem.MinValue
	}
	if metadataItem.MaxValue != nil {
		metadataMap["max_value"] = metadataItem.MaxValue
	}
	if metadataItem.MinLength != nil {
		metadataMap["min_length"] = metadataItem.MinLength
	}
	if metadataItem.MaxLength != nil {
		metadataMap["max_length"] = metadataItem.MaxLength
	}
	if metadataItem.Matches != nil {
		metadataMap["matches"] = metadataItem.Matches
	}
	if metadataItem.Position != nil {
		metadataMap["position"] = metadataItem.Position
	}
	if metadataItem.GroupBy != nil {
		metadataMap["group_by"] = metadataItem.GroupBy
	}
	if metadataItem.Source != nil {
		metadataMap["source"] = metadataItem.Source
	}

	return metadataMap
}

func dataSourceActionFlattenOutputs(result []schematicsv1.VariableData) (actionOutputs []map[string]interface{}) {
	for _, actionOutputsItem := range result {
		actionOutputs = append(actionOutputs, dataSourceActionOutputsToMap(actionOutputsItem))
	}

	return actionOutputs
}

func dataSourceActionOutputsToMap(outputsItem schematicsv1.VariableData) (outputsMap map[string]interface{}) {
	outputsMap = map[string]interface{}{}

	if outputsItem.Name != nil {
		outputsMap["name"] = outputsItem.Name
	}
	if outputsItem.Value != nil {
		outputsMap["value"] = outputsItem.Value
	}
	if outputsItem.Metadata != nil {
		metadataList := []map[string]interface{}{}
		metadataMap := dataSourceActionOutputsMetadataToMap(*outputsItem.Metadata)
		metadataList = append(metadataList, metadataMap)
		outputsMap["metadata"] = metadataList
	}
	if outputsItem.Link != nil {
		outputsMap["link"] = outputsItem.Link
	}

	return outputsMap
}

func dataSourceActionOutputsMetadataToMap(metadataItem schematicsv1.VariableMetadata) (metadataMap map[string]interface{}) {
	metadataMap = map[string]interface{}{}

	if metadataItem.Type != nil {
		metadataMap["type"] = metadataItem.Type
	}
	if metadataItem.Aliases != nil {
		metadataMap["aliases"] = metadataItem.Aliases
	}
	if metadataItem.Description != nil {
		metadataMap["description"] = metadataItem.Description
	}
	if metadataItem.DefaultValue != nil {
		metadataMap["default_value"] = metadataItem.DefaultValue
	}
	if metadataItem.Secure != nil {
		metadataMap["secure"] = metadataItem.Secure
	}
	if metadataItem.Immutable != nil {
		metadataMap["immutable"] = metadataItem.Immutable
	}
	if metadataItem.Hidden != nil {
		metadataMap["hidden"] = metadataItem.Hidden
	}
	if metadataItem.Options != nil {
		metadataMap["options"] = metadataItem.Options
	}
	if metadataItem.MinValue != nil {
		metadataMap["min_value"] = metadataItem.MinValue
	}
	if metadataItem.MaxValue != nil {
		metadataMap["max_value"] = metadataItem.MaxValue
	}
	if metadataItem.MinLength != nil {
		metadataMap["min_length"] = metadataItem.MinLength
	}
	if metadataItem.MaxLength != nil {
		metadataMap["max_length"] = metadataItem.MaxLength
	}
	if metadataItem.Matches != nil {
		metadataMap["matches"] = metadataItem.Matches
	}
	if metadataItem.Position != nil {
		metadataMap["position"] = metadataItem.Position
	}
	if metadataItem.GroupBy != nil {
		metadataMap["group_by"] = metadataItem.GroupBy
	}
	if metadataItem.Source != nil {
		metadataMap["source"] = metadataItem.Source
	}

	return metadataMap
}

func dataSourceActionFlattenSettings(result []schematicsv1.VariableData) (settings []map[string]interface{}) {
	for _, settingsItem := range result {
		settings = append(settings, dataSourceActionSettingsToMap(settingsItem))
	}

	return settings
}

func dataSourceActionSettingsToMap(settingsItem schematicsv1.VariableData) (settingsMap map[string]interface{}) {
	settingsMap = map[string]interface{}{}

	if settingsItem.Name != nil {
		settingsMap["name"] = settingsItem.Name
	}
	if settingsItem.Value != nil {
		settingsMap["value"] = settingsItem.Value
	}
	if settingsItem.Metadata != nil {
		metadataList := []map[string]interface{}{}
		metadataMap := dataSourceActionSettingsMetadataToMap(*settingsItem.Metadata)
		metadataList = append(metadataList, metadataMap)
		settingsMap["metadata"] = metadataList
	}
	if settingsItem.Link != nil {
		settingsMap["link"] = settingsItem.Link
	}

	return settingsMap
}

func dataSourceActionSettingsMetadataToMap(metadataItem schematicsv1.VariableMetadata) (metadataMap map[string]interface{}) {
	metadataMap = map[string]interface{}{}

	if metadataItem.Type != nil {
		metadataMap["type"] = metadataItem.Type
	}
	if metadataItem.Aliases != nil {
		metadataMap["aliases"] = metadataItem.Aliases
	}
	if metadataItem.Description != nil {
		metadataMap["description"] = metadataItem.Description
	}
	if metadataItem.DefaultValue != nil {
		metadataMap["default_value"] = metadataItem.DefaultValue
	}
	if metadataItem.Secure != nil {
		metadataMap["secure"] = metadataItem.Secure
	}
	if metadataItem.Immutable != nil {
		metadataMap["immutable"] = metadataItem.Immutable
	}
	if metadataItem.Hidden != nil {
		metadataMap["hidden"] = metadataItem.Hidden
	}
	if metadataItem.Options != nil {
		metadataMap["options"] = metadataItem.Options
	}
	if metadataItem.MinValue != nil {
		metadataMap["min_value"] = metadataItem.MinValue
	}
	if metadataItem.MaxValue != nil {
		metadataMap["max_value"] = metadataItem.MaxValue
	}
	if metadataItem.MinLength != nil {
		metadataMap["min_length"] = metadataItem.MinLength
	}
	if metadataItem.MaxLength != nil {
		metadataMap["max_length"] = metadataItem.MaxLength
	}
	if metadataItem.Matches != nil {
		metadataMap["matches"] = metadataItem.Matches
	}
	if metadataItem.Position != nil {
		metadataMap["position"] = metadataItem.Position
	}
	if metadataItem.GroupBy != nil {
		metadataMap["group_by"] = metadataItem.GroupBy
	}
	if metadataItem.Source != nil {
		metadataMap["source"] = metadataItem.Source
	}

	return metadataMap
}

func dataSourceActionFlattenState(result schematicsv1.ActionState) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceActionStateToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceActionStateToMap(stateItem schematicsv1.ActionState) (stateMap map[string]interface{}) {
	stateMap = map[string]interface{}{}

	if stateItem.StatusCode != nil {
		stateMap["status_code"] = stateItem.StatusCode
	}
	if stateItem.StatusJobID != nil {
		stateMap["status_job_id"] = stateItem.StatusJobID
	}
	if stateItem.StatusMessage != nil {
		stateMap["status_message"] = stateItem.StatusMessage
	}

	return stateMap
}

func dataSourceActionFlattenSysLock(result schematicsv1.SystemLock) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceActionSysLockToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceActionSysLockToMap(sysLockItem schematicsv1.SystemLock) (sysLockMap map[string]interface{}) {
	sysLockMap = map[string]interface{}{}

	if sysLockItem.SysLocked != nil {
		sysLockMap["sys_locked"] = sysLockItem.SysLocked
	}
	if sysLockItem.SysLockedBy != nil {
		sysLockMap["sys_locked_by"] = sysLockItem.SysLockedBy
	}
	if sysLockItem.SysLockedAt != nil {
		sysLockMap["sys_locked_at"] = sysLockItem.SysLockedAt.String()
	}

	return sysLockMap
}
