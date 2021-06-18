package bsdata

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/protocol/packp/sideband"
)

const (
	baseDataRepoURL = "https://github.com/BSData"
	directory       = "./checkout-tmp"
)

type Catalogue struct {
	XMLName             xml.Name `xml:"catalogue"`
	Text                string   `xml:",chardata"`
	ID                  string   `xml:"id,attr"`
	Name                string   `xml:"name,attr"`
	Revision            string   `xml:"revision,attr"`
	BattleScribeVersion string   `xml:"battleScribeVersion,attr"`
	AuthorName          string   `xml:"authorName,attr"`
	AuthorContact       string   `xml:"authorContact,attr"`
	AuthorUrl           string   `xml:"authorUrl,attr"`
	Library             string   `xml:"library,attr"`
	GameSystemId        string   `xml:"gameSystemId,attr"`
	GameSystemRevision  string   `xml:"gameSystemRevision,attr"`
	Xmlns               string   `xml:"xmlns,attr"`
	Comment             string   `xml:"comment"`
	Readme              string   `xml:"readme"`
	Publications        struct {
		Text        string `xml:",chardata"`
		Publication []struct {
			Text            string `xml:",chardata"`
			ID              string `xml:"id,attr"`
			Name            string `xml:"name,attr"`
			ShortName       string `xml:"shortName,attr"`
			Publisher       string `xml:"publisher,attr"`
			PublicationDate string `xml:"publicationDate,attr"`
		} `xml:"publication"`
	} `xml:"publications"`
	ProfileTypes struct {
		Text        string `xml:",chardata"`
		ProfileType struct {
			Text                string `xml:",chardata"`
			ID                  string `xml:"id,attr"`
			Name                string `xml:"name,attr"`
			CharacteristicTypes struct {
				Text               string `xml:",chardata"`
				CharacteristicType []struct {
					Text string `xml:",chardata"`
					ID   string `xml:"id,attr"`
					Name string `xml:"name,attr"`
				} `xml:"characteristicType"`
			} `xml:"characteristicTypes"`
		} `xml:"profileType"`
	} `xml:"profileTypes"`
	CategoryEntries struct {
		Text          string `xml:",chardata"`
		CategoryEntry []struct {
			Text   string `xml:",chardata"`
			ID     string `xml:"id,attr"`
			Name   string `xml:"name,attr"`
			Hidden string `xml:"hidden,attr"`
		} `xml:"categoryEntry"`
	} `xml:"categoryEntries"`
	EntryLinks struct {
		Text      string `xml:",chardata"`
		EntryLink []struct {
			Text       string `xml:",chardata"`
			ID         string `xml:"id,attr"`
			Name       string `xml:"name,attr"`
			Hidden     string `xml:"hidden,attr"`
			Collective string `xml:"collective,attr"`
			Import     string `xml:"import,attr"`
			TargetId   string `xml:"targetId,attr"`
			Type       string `xml:"type,attr"`
			Modifiers  struct {
				Text     string `xml:",chardata"`
				Modifier struct {
					Text       string `xml:",chardata"`
					Type       string `xml:"type,attr"`
					Field      string `xml:"field,attr"`
					Value      string `xml:"value,attr"`
					Conditions struct {
						Text      string `xml:",chardata"`
						Condition struct {
							Text                   string `xml:",chardata"`
							Field                  string `xml:"field,attr"`
							Scope                  string `xml:"scope,attr"`
							Value                  string `xml:"value,attr"`
							PercentValue           string `xml:"percentValue,attr"`
							Shared                 string `xml:"shared,attr"`
							IncludeChildSelections string `xml:"includeChildSelections,attr"`
							IncludeChildForces     string `xml:"includeChildForces,attr"`
							ChildId                string `xml:"childId,attr"`
							Type                   string `xml:"type,attr"`
						} `xml:"condition"`
					} `xml:"conditions"`
				} `xml:"modifier"`
			} `xml:"modifiers"`
		} `xml:"entryLink"`
	} `xml:"entryLinks"`
	SharedSelectionEntries struct {
		Text           string `xml:",chardata"`
		SelectionEntry []struct {
			Text          string `xml:",chardata"`
			ID            string `xml:"id,attr"`
			Name          string `xml:"name,attr"`
			PublicationId string `xml:"publicationId,attr"`
			Page          string `xml:"page,attr"`
			Hidden        string `xml:"hidden,attr"`
			Collective    string `xml:"collective,attr"`
			Import        string `xml:"import,attr"`
			Type          string `xml:"type,attr"`
			Constraints   struct {
				Text       string `xml:",chardata"`
				Constraint []struct {
					Text                   string `xml:",chardata"`
					Field                  string `xml:"field,attr"`
					Scope                  string `xml:"scope,attr"`
					Value                  string `xml:"value,attr"`
					PercentValue           string `xml:"percentValue,attr"`
					Shared                 string `xml:"shared,attr"`
					IncludeChildSelections string `xml:"includeChildSelections,attr"`
					IncludeChildForces     string `xml:"includeChildForces,attr"`
					ID                     string `xml:"id,attr"`
					Type                   string `xml:"type,attr"`
				} `xml:"constraint"`
			} `xml:"constraints"`
			Profiles struct {
				Text    string `xml:",chardata"`
				Profile []struct {
					Text            string `xml:",chardata"`
					ID              string `xml:"id,attr"`
					Name            string `xml:"name,attr"`
					Hidden          string `xml:"hidden,attr"`
					TypeId          string `xml:"typeId,attr"`
					TypeName        string `xml:"typeName,attr"`
					PublicationId   string `xml:"publicationId,attr"`
					Page            string `xml:"page,attr"`
					Characteristics struct {
						Text           string `xml:",chardata"`
						Characteristic []struct {
							Text   string `xml:",chardata"`
							Name   string `xml:"name,attr"`
							TypeId string `xml:"typeId,attr"`
						} `xml:"characteristic"`
					} `xml:"characteristics"`
				} `xml:"profile"`
			} `xml:"profiles"`
			InfoLinks struct {
				Text     string `xml:",chardata"`
				InfoLink []struct {
					Text      string `xml:",chardata"`
					ID        string `xml:"id,attr"`
					Name      string `xml:"name,attr"`
					Hidden    string `xml:"hidden,attr"`
					TargetId  string `xml:"targetId,attr"`
					Type      string `xml:"type,attr"`
					Modifiers struct {
						Text     string `xml:",chardata"`
						Modifier struct {
							Text       string `xml:",chardata"`
							Type       string `xml:"type,attr"`
							Field      string `xml:"field,attr"`
							Value      string `xml:"value,attr"`
							Conditions struct {
								Text      string `xml:",chardata"`
								Condition struct {
									Text                   string `xml:",chardata"`
									Field                  string `xml:"field,attr"`
									Scope                  string `xml:"scope,attr"`
									Value                  string `xml:"value,attr"`
									PercentValue           string `xml:"percentValue,attr"`
									Shared                 string `xml:"shared,attr"`
									IncludeChildSelections string `xml:"includeChildSelections,attr"`
									IncludeChildForces     string `xml:"includeChildForces,attr"`
									ChildId                string `xml:"childId,attr"`
									Type                   string `xml:"type,attr"`
								} `xml:"condition"`
							} `xml:"conditions"`
						} `xml:"modifier"`
					} `xml:"modifiers"`
				} `xml:"infoLink"`
			} `xml:"infoLinks"`
			CategoryLinks struct {
				Text         string `xml:",chardata"`
				CategoryLink []struct {
					Text     string `xml:",chardata"`
					ID       string `xml:"id,attr"`
					Name     string `xml:"name,attr"`
					Hidden   string `xml:"hidden,attr"`
					TargetId string `xml:"targetId,attr"`
					Primary  string `xml:"primary,attr"`
				} `xml:"categoryLink"`
			} `xml:"categoryLinks"`
			SelectionEntries struct {
				Text           string `xml:",chardata"`
				SelectionEntry []struct {
					Text        string `xml:",chardata"`
					ID          string `xml:"id,attr"`
					Name        string `xml:"name,attr"`
					Hidden      string `xml:"hidden,attr"`
					Collective  string `xml:"collective,attr"`
					Import      string `xml:"import,attr"`
					Type        string `xml:"type,attr"`
					Constraints struct {
						Text       string `xml:",chardata"`
						Constraint []struct {
							Text                   string `xml:",chardata"`
							Field                  string `xml:"field,attr"`
							Scope                  string `xml:"scope,attr"`
							Value                  string `xml:"value,attr"`
							PercentValue           string `xml:"percentValue,attr"`
							Shared                 string `xml:"shared,attr"`
							IncludeChildSelections string `xml:"includeChildSelections,attr"`
							IncludeChildForces     string `xml:"includeChildForces,attr"`
							ID                     string `xml:"id,attr"`
							Type                   string `xml:"type,attr"`
						} `xml:"constraint"`
					} `xml:"constraints"`
					Profiles struct {
						Text    string `xml:",chardata"`
						Profile struct {
							Text            string `xml:",chardata"`
							ID              string `xml:"id,attr"`
							Name            string `xml:"name,attr"`
							Hidden          string `xml:"hidden,attr"`
							TypeId          string `xml:"typeId,attr"`
							TypeName        string `xml:"typeName,attr"`
							Characteristics struct {
								Text           string `xml:",chardata"`
								Characteristic []struct {
									Text   string `xml:",chardata"`
									Name   string `xml:"name,attr"`
									TypeId string `xml:"typeId,attr"`
								} `xml:"characteristic"`
							} `xml:"characteristics"`
						} `xml:"profile"`
					} `xml:"profiles"`
					Costs struct {
						Text string `xml:",chardata"`
						Cost []struct {
							Text   string `xml:",chardata"`
							Name   string `xml:"name,attr"`
							TypeId string `xml:"typeId,attr"`
							Value  string `xml:"value,attr"`
						} `xml:"cost"`
					} `xml:"costs"`
				} `xml:"selectionEntry"`
			} `xml:"selectionEntries"`
			EntryLinks struct {
				Text      string `xml:",chardata"`
				EntryLink []struct {
					Text        string `xml:",chardata"`
					ID          string `xml:"id,attr"`
					Name        string `xml:"name,attr"`
					Hidden      string `xml:"hidden,attr"`
					Collective  string `xml:"collective,attr"`
					Import      string `xml:"import,attr"`
					TargetId    string `xml:"targetId,attr"`
					Type        string `xml:"type,attr"`
					Constraints struct {
						Text       string `xml:",chardata"`
						Constraint struct {
							Text                   string `xml:",chardata"`
							Field                  string `xml:"field,attr"`
							Scope                  string `xml:"scope,attr"`
							Value                  string `xml:"value,attr"`
							PercentValue           string `xml:"percentValue,attr"`
							Shared                 string `xml:"shared,attr"`
							IncludeChildSelections string `xml:"includeChildSelections,attr"`
							IncludeChildForces     string `xml:"includeChildForces,attr"`
							ID                     string `xml:"id,attr"`
							Type                   string `xml:"type,attr"`
						} `xml:"constraint"`
					} `xml:"constraints"`
					Modifiers struct {
						Text     string `xml:",chardata"`
						Modifier []struct {
							Text            string `xml:",chardata"`
							Type            string `xml:"type,attr"`
							Field           string `xml:"field,attr"`
							Value           string `xml:"value,attr"`
							ConditionGroups struct {
								Text           string `xml:",chardata"`
								ConditionGroup struct {
									Text       string `xml:",chardata"`
									Type       string `xml:"type,attr"`
									Conditions struct {
										Text      string `xml:",chardata"`
										Condition []struct {
											Text                   string `xml:",chardata"`
											Field                  string `xml:"field,attr"`
											Scope                  string `xml:"scope,attr"`
											Value                  string `xml:"value,attr"`
											PercentValue           string `xml:"percentValue,attr"`
											Shared                 string `xml:"shared,attr"`
											IncludeChildSelections string `xml:"includeChildSelections,attr"`
											IncludeChildForces     string `xml:"includeChildForces,attr"`
											ChildId                string `xml:"childId,attr"`
											Type                   string `xml:"type,attr"`
										} `xml:"condition"`
									} `xml:"conditions"`
								} `xml:"conditionGroup"`
							} `xml:"conditionGroups"`
						} `xml:"modifier"`
					} `xml:"modifiers"`
				} `xml:"entryLink"`
			} `xml:"entryLinks"`
			Costs struct {
				Text string `xml:",chardata"`
				Cost []struct {
					Text   string `xml:",chardata"`
					Name   string `xml:"name,attr"`
					TypeId string `xml:"typeId,attr"`
					Value  string `xml:"value,attr"`
				} `xml:"cost"`
			} `xml:"costs"`
			Modifiers struct {
				Text     string `xml:",chardata"`
				Modifier []struct {
					Text       string `xml:",chardata"`
					Type       string `xml:"type,attr"`
					Field      string `xml:"field,attr"`
					Value      string `xml:"value,attr"`
					Conditions struct {
						Text      string `xml:",chardata"`
						Condition struct {
							Text                   string `xml:",chardata"`
							Field                  string `xml:"field,attr"`
							Scope                  string `xml:"scope,attr"`
							Value                  string `xml:"value,attr"`
							PercentValue           string `xml:"percentValue,attr"`
							Shared                 string `xml:"shared,attr"`
							IncludeChildSelections string `xml:"includeChildSelections,attr"`
							IncludeChildForces     string `xml:"includeChildForces,attr"`
							ChildId                string `xml:"childId,attr"`
							Type                   string `xml:"type,attr"`
						} `xml:"condition"`
					} `xml:"conditions"`
					ConditionGroups struct {
						Text           string `xml:",chardata"`
						ConditionGroup struct {
							Text       string `xml:",chardata"`
							Type       string `xml:"type,attr"`
							Conditions struct {
								Text      string `xml:",chardata"`
								Condition []struct {
									Text                   string `xml:",chardata"`
									Field                  string `xml:"field,attr"`
									Scope                  string `xml:"scope,attr"`
									Value                  string `xml:"value,attr"`
									PercentValue           string `xml:"percentValue,attr"`
									Shared                 string `xml:"shared,attr"`
									IncludeChildSelections string `xml:"includeChildSelections,attr"`
									IncludeChildForces     string `xml:"includeChildForces,attr"`
									ChildId                string `xml:"childId,attr"`
									Type                   string `xml:"type,attr"`
								} `xml:"condition"`
							} `xml:"conditions"`
						} `xml:"conditionGroup"`
					} `xml:"conditionGroups"`
					Repeats struct {
						Text   string `xml:",chardata"`
						Repeat struct {
							Text                   string `xml:",chardata"`
							Field                  string `xml:"field,attr"`
							Scope                  string `xml:"scope,attr"`
							Value                  string `xml:"value,attr"`
							PercentValue           string `xml:"percentValue,attr"`
							Shared                 string `xml:"shared,attr"`
							IncludeChildSelections string `xml:"includeChildSelections,attr"`
							IncludeChildForces     string `xml:"includeChildForces,attr"`
							ChildId                string `xml:"childId,attr"`
							Repeats                string `xml:"repeats,attr"`
							RoundUp                string `xml:"roundUp,attr"`
						} `xml:"repeat"`
					} `xml:"repeats"`
				} `xml:"modifier"`
			} `xml:"modifiers"`
		} `xml:"selectionEntry"`
	} `xml:"sharedSelectionEntries"`
	SharedSelectionEntryGroups struct {
		Text                string `xml:",chardata"`
		SelectionEntryGroup []struct {
			Text          string `xml:",chardata"`
			ID            string `xml:"id,attr"`
			Name          string `xml:"name,attr"`
			PublicationId string `xml:"publicationId,attr"`
			Hidden        string `xml:"hidden,attr"`
			Collective    string `xml:"collective,attr"`
			Import        string `xml:"import,attr"`
			Modifiers     struct {
				Text     string `xml:",chardata"`
				Modifier []struct {
					Text            string `xml:",chardata"`
					Type            string `xml:"type,attr"`
					Field           string `xml:"field,attr"`
					Value           string `xml:"value,attr"`
					ConditionGroups struct {
						Text           string `xml:",chardata"`
						ConditionGroup struct {
							Text       string `xml:",chardata"`
							Type       string `xml:"type,attr"`
							Conditions struct {
								Text      string `xml:",chardata"`
								Condition []struct {
									Text                   string `xml:",chardata"`
									Field                  string `xml:"field,attr"`
									Scope                  string `xml:"scope,attr"`
									Value                  string `xml:"value,attr"`
									PercentValue           string `xml:"percentValue,attr"`
									Shared                 string `xml:"shared,attr"`
									IncludeChildSelections string `xml:"includeChildSelections,attr"`
									IncludeChildForces     string `xml:"includeChildForces,attr"`
									ChildId                string `xml:"childId,attr"`
									Type                   string `xml:"type,attr"`
								} `xml:"condition"`
							} `xml:"conditions"`
						} `xml:"conditionGroup"`
					} `xml:"conditionGroups"`
					Conditions struct {
						Text      string `xml:",chardata"`
						Condition struct {
							Text                   string `xml:",chardata"`
							Field                  string `xml:"field,attr"`
							Scope                  string `xml:"scope,attr"`
							Value                  string `xml:"value,attr"`
							PercentValue           string `xml:"percentValue,attr"`
							Shared                 string `xml:"shared,attr"`
							IncludeChildSelections string `xml:"includeChildSelections,attr"`
							IncludeChildForces     string `xml:"includeChildForces,attr"`
							ChildId                string `xml:"childId,attr"`
							Type                   string `xml:"type,attr"`
						} `xml:"condition"`
					} `xml:"conditions"`
					Repeats struct {
						Text   string `xml:",chardata"`
						Repeat struct {
							Text                   string `xml:",chardata"`
							Field                  string `xml:"field,attr"`
							Scope                  string `xml:"scope,attr"`
							Value                  string `xml:"value,attr"`
							PercentValue           string `xml:"percentValue,attr"`
							Shared                 string `xml:"shared,attr"`
							IncludeChildSelections string `xml:"includeChildSelections,attr"`
							IncludeChildForces     string `xml:"includeChildForces,attr"`
							ChildId                string `xml:"childId,attr"`
							Repeats                string `xml:"repeats,attr"`
							RoundUp                string `xml:"roundUp,attr"`
						} `xml:"repeat"`
					} `xml:"repeats"`
				} `xml:"modifier"`
			} `xml:"modifiers"`
			SelectionEntries struct {
				Text           string `xml:",chardata"`
				SelectionEntry []struct {
					Text          string `xml:",chardata"`
					ID            string `xml:"id,attr"`
					Name          string `xml:"name,attr"`
					Hidden        string `xml:"hidden,attr"`
					Collective    string `xml:"collective,attr"`
					Import        string `xml:"import,attr"`
					Type          string `xml:"type,attr"`
					PublicationId string `xml:"publicationId,attr"`
					Page          string `xml:"page,attr"`
					Constraints   struct {
						Text       string `xml:",chardata"`
						Constraint []struct {
							Text                   string `xml:",chardata"`
							Field                  string `xml:"field,attr"`
							Scope                  string `xml:"scope,attr"`
							Value                  string `xml:"value,attr"`
							PercentValue           string `xml:"percentValue,attr"`
							Shared                 string `xml:"shared,attr"`
							IncludeChildSelections string `xml:"includeChildSelections,attr"`
							IncludeChildForces     string `xml:"includeChildForces,attr"`
							ID                     string `xml:"id,attr"`
							Type                   string `xml:"type,attr"`
						} `xml:"constraint"`
					} `xml:"constraints"`
					Profiles struct {
						Text    string `xml:",chardata"`
						Profile struct {
							Text            string `xml:",chardata"`
							ID              string `xml:"id,attr"`
							Name            string `xml:"name,attr"`
							Hidden          string `xml:"hidden,attr"`
							TypeId          string `xml:"typeId,attr"`
							TypeName        string `xml:"typeName,attr"`
							Characteristics struct {
								Text           string `xml:",chardata"`
								Characteristic []struct {
									Text   string `xml:",chardata"`
									Name   string `xml:"name,attr"`
									TypeId string `xml:"typeId,attr"`
								} `xml:"characteristic"`
							} `xml:"characteristics"`
						} `xml:"profile"`
					} `xml:"profiles"`
					Costs struct {
						Text string `xml:",chardata"`
						Cost []struct {
							Text   string `xml:",chardata"`
							Name   string `xml:"name,attr"`
							TypeId string `xml:"typeId,attr"`
							Value  string `xml:"value,attr"`
						} `xml:"cost"`
					} `xml:"costs"`
					SelectionEntryGroups struct {
						Text                string `xml:",chardata"`
						SelectionEntryGroup struct {
							Text           string `xml:",chardata"`
							ID             string `xml:"id,attr"`
							Name           string `xml:"name,attr"`
							Hidden         string `xml:"hidden,attr"`
							Collective     string `xml:"collective,attr"`
							Import         string `xml:"import,attr"`
							ModifierGroups struct {
								Text          string `xml:",chardata"`
								ModifierGroup struct {
									Text       string `xml:",chardata"`
									Conditions struct {
										Text      string `xml:",chardata"`
										Condition struct {
											Text                   string `xml:",chardata"`
											Field                  string `xml:"field,attr"`
											Scope                  string `xml:"scope,attr"`
											Value                  string `xml:"value,attr"`
											PercentValue           string `xml:"percentValue,attr"`
											Shared                 string `xml:"shared,attr"`
											IncludeChildSelections string `xml:"includeChildSelections,attr"`
											IncludeChildForces     string `xml:"includeChildForces,attr"`
											ChildId                string `xml:"childId,attr"`
											Type                   string `xml:"type,attr"`
										} `xml:"condition"`
									} `xml:"conditions"`
									Modifiers struct {
										Text     string `xml:",chardata"`
										Modifier []struct {
											Text  string `xml:",chardata"`
											Type  string `xml:"type,attr"`
											Field string `xml:"field,attr"`
											Value string `xml:"value,attr"`
										} `xml:"modifier"`
									} `xml:"modifiers"`
									ConditionGroups struct {
										Text           string `xml:",chardata"`
										ConditionGroup struct {
											Text       string `xml:",chardata"`
											Type       string `xml:"type,attr"`
											Conditions struct {
												Text      string `xml:",chardata"`
												Condition []struct {
													Text                   string `xml:",chardata"`
													Field                  string `xml:"field,attr"`
													Scope                  string `xml:"scope,attr"`
													Value                  string `xml:"value,attr"`
													PercentValue           string `xml:"percentValue,attr"`
													Shared                 string `xml:"shared,attr"`
													IncludeChildSelections string `xml:"includeChildSelections,attr"`
													IncludeChildForces     string `xml:"includeChildForces,attr"`
													ChildId                string `xml:"childId,attr"`
													Type                   string `xml:"type,attr"`
												} `xml:"condition"`
											} `xml:"conditions"`
										} `xml:"conditionGroup"`
									} `xml:"conditionGroups"`
								} `xml:"modifierGroup"`
							} `xml:"modifierGroups"`
							Constraints struct {
								Text       string `xml:",chardata"`
								Constraint []struct {
									Text                   string `xml:",chardata"`
									Field                  string `xml:"field,attr"`
									Scope                  string `xml:"scope,attr"`
									Value                  string `xml:"value,attr"`
									PercentValue           string `xml:"percentValue,attr"`
									Shared                 string `xml:"shared,attr"`
									IncludeChildSelections string `xml:"includeChildSelections,attr"`
									IncludeChildForces     string `xml:"includeChildForces,attr"`
									ID                     string `xml:"id,attr"`
									Type                   string `xml:"type,attr"`
								} `xml:"constraint"`
							} `xml:"constraints"`
							SelectionEntries struct {
								Text           string `xml:",chardata"`
								SelectionEntry []struct {
									Text        string `xml:",chardata"`
									ID          string `xml:"id,attr"`
									Name        string `xml:"name,attr"`
									Hidden      string `xml:"hidden,attr"`
									Collective  string `xml:"collective,attr"`
									Import      string `xml:"import,attr"`
									Type        string `xml:"type,attr"`
									Constraints struct {
										Text       string `xml:",chardata"`
										Constraint struct {
											Text                   string `xml:",chardata"`
											Field                  string `xml:"field,attr"`
											Scope                  string `xml:"scope,attr"`
											Value                  string `xml:"value,attr"`
											PercentValue           string `xml:"percentValue,attr"`
											Shared                 string `xml:"shared,attr"`
											IncludeChildSelections string `xml:"includeChildSelections,attr"`
											IncludeChildForces     string `xml:"includeChildForces,attr"`
											ID                     string `xml:"id,attr"`
											Type                   string `xml:"type,attr"`
										} `xml:"constraint"`
									} `xml:"constraints"`
									Profiles struct {
										Text    string `xml:",chardata"`
										Profile struct {
											Text            string `xml:",chardata"`
											ID              string `xml:"id,attr"`
											Name            string `xml:"name,attr"`
											Hidden          string `xml:"hidden,attr"`
											TypeId          string `xml:"typeId,attr"`
											TypeName        string `xml:"typeName,attr"`
											Characteristics struct {
												Text           string `xml:",chardata"`
												Characteristic struct {
													Text   string `xml:",chardata"`
													Name   string `xml:"name,attr"`
													TypeId string `xml:"typeId,attr"`
												} `xml:"characteristic"`
											} `xml:"characteristics"`
										} `xml:"profile"`
									} `xml:"profiles"`
									Costs struct {
										Text string `xml:",chardata"`
										Cost []struct {
											Text   string `xml:",chardata"`
											Name   string `xml:"name,attr"`
											TypeId string `xml:"typeId,attr"`
											Value  string `xml:"value,attr"`
										} `xml:"cost"`
									} `xml:"costs"`
								} `xml:"selectionEntry"`
							} `xml:"selectionEntries"`
						} `xml:"selectionEntryGroup"`
					} `xml:"selectionEntryGroups"`
					Modifiers struct {
						Text     string `xml:",chardata"`
						Modifier struct {
							Text            string `xml:",chardata"`
							Type            string `xml:"type,attr"`
							Field           string `xml:"field,attr"`
							Value           string `xml:"value,attr"`
							ConditionGroups struct {
								Text           string `xml:",chardata"`
								ConditionGroup struct {
									Text       string `xml:",chardata"`
									Type       string `xml:"type,attr"`
									Conditions struct {
										Text      string `xml:",chardata"`
										Condition []struct {
											Text                   string `xml:",chardata"`
											Field                  string `xml:"field,attr"`
											Scope                  string `xml:"scope,attr"`
											Value                  string `xml:"value,attr"`
											PercentValue           string `xml:"percentValue,attr"`
											Shared                 string `xml:"shared,attr"`
											IncludeChildSelections string `xml:"includeChildSelections,attr"`
											IncludeChildForces     string `xml:"includeChildForces,attr"`
											ChildId                string `xml:"childId,attr"`
											Type                   string `xml:"type,attr"`
										} `xml:"condition"`
									} `xml:"conditions"`
								} `xml:"conditionGroup"`
							} `xml:"conditionGroups"`
							Conditions struct {
								Text      string `xml:",chardata"`
								Condition struct {
									Text                   string `xml:",chardata"`
									Field                  string `xml:"field,attr"`
									Scope                  string `xml:"scope,attr"`
									Value                  string `xml:"value,attr"`
									PercentValue           string `xml:"percentValue,attr"`
									Shared                 string `xml:"shared,attr"`
									IncludeChildSelections string `xml:"includeChildSelections,attr"`
									IncludeChildForces     string `xml:"includeChildForces,attr"`
									ChildId                string `xml:"childId,attr"`
									Type                   string `xml:"type,attr"`
								} `xml:"condition"`
							} `xml:"conditions"`
						} `xml:"modifier"`
					} `xml:"modifiers"`
					CategoryLinks struct {
						Text         string `xml:",chardata"`
						CategoryLink struct {
							Text     string `xml:",chardata"`
							ID       string `xml:"id,attr"`
							Name     string `xml:"name,attr"`
							Hidden   string `xml:"hidden,attr"`
							TargetId string `xml:"targetId,attr"`
							Primary  string `xml:"primary,attr"`
						} `xml:"categoryLink"`
					} `xml:"categoryLinks"`
				} `xml:"selectionEntry"`
			} `xml:"selectionEntries"`
			Constraints struct {
				Text       string `xml:",chardata"`
				Constraint []struct {
					Text                   string `xml:",chardata"`
					Field                  string `xml:"field,attr"`
					Scope                  string `xml:"scope,attr"`
					Value                  string `xml:"value,attr"`
					PercentValue           string `xml:"percentValue,attr"`
					Shared                 string `xml:"shared,attr"`
					IncludeChildSelections string `xml:"includeChildSelections,attr"`
					IncludeChildForces     string `xml:"includeChildForces,attr"`
					ID                     string `xml:"id,attr"`
					Type                   string `xml:"type,attr"`
				} `xml:"constraint"`
			} `xml:"constraints"`
			EntryLinks struct {
				Text      string `xml:",chardata"`
				EntryLink []struct {
					Text       string `xml:",chardata"`
					ID         string `xml:"id,attr"`
					Name       string `xml:"name,attr"`
					Hidden     string `xml:"hidden,attr"`
					Collective string `xml:"collective,attr"`
					Import     string `xml:"import,attr"`
					TargetId   string `xml:"targetId,attr"`
					Type       string `xml:"type,attr"`
					Modifiers  struct {
						Text     string `xml:",chardata"`
						Modifier struct {
							Text            string `xml:",chardata"`
							Type            string `xml:"type,attr"`
							Field           string `xml:"field,attr"`
							Value           string `xml:"value,attr"`
							ConditionGroups struct {
								Text           string `xml:",chardata"`
								ConditionGroup struct {
									Text       string `xml:",chardata"`
									Type       string `xml:"type,attr"`
									Conditions struct {
										Text      string `xml:",chardata"`
										Condition []struct {
											Text                   string `xml:",chardata"`
											Field                  string `xml:"field,attr"`
											Scope                  string `xml:"scope,attr"`
											Value                  string `xml:"value,attr"`
											PercentValue           string `xml:"percentValue,attr"`
											Shared                 string `xml:"shared,attr"`
											IncludeChildSelections string `xml:"includeChildSelections,attr"`
											IncludeChildForces     string `xml:"includeChildForces,attr"`
											ChildId                string `xml:"childId,attr"`
											Type                   string `xml:"type,attr"`
										} `xml:"condition"`
									} `xml:"conditions"`
									ConditionGroups struct {
										Text           string `xml:",chardata"`
										ConditionGroup struct {
											Text       string `xml:",chardata"`
											Type       string `xml:"type,attr"`
											Conditions struct {
												Text      string `xml:",chardata"`
												Condition []struct {
													Text                   string `xml:",chardata"`
													Field                  string `xml:"field,attr"`
													Scope                  string `xml:"scope,attr"`
													Value                  string `xml:"value,attr"`
													PercentValue           string `xml:"percentValue,attr"`
													Shared                 string `xml:"shared,attr"`
													IncludeChildSelections string `xml:"includeChildSelections,attr"`
													IncludeChildForces     string `xml:"includeChildForces,attr"`
													ChildId                string `xml:"childId,attr"`
													Type                   string `xml:"type,attr"`
												} `xml:"condition"`
											} `xml:"conditions"`
										} `xml:"conditionGroup"`
									} `xml:"conditionGroups"`
								} `xml:"conditionGroup"`
							} `xml:"conditionGroups"`
							Conditions struct {
								Text      string `xml:",chardata"`
								Condition struct {
									Text                   string `xml:",chardata"`
									Field                  string `xml:"field,attr"`
									Scope                  string `xml:"scope,attr"`
									Value                  string `xml:"value,attr"`
									PercentValue           string `xml:"percentValue,attr"`
									Shared                 string `xml:"shared,attr"`
									IncludeChildSelections string `xml:"includeChildSelections,attr"`
									IncludeChildForces     string `xml:"includeChildForces,attr"`
									ChildId                string `xml:"childId,attr"`
									Type                   string `xml:"type,attr"`
								} `xml:"condition"`
							} `xml:"conditions"`
						} `xml:"modifier"`
					} `xml:"modifiers"`
					Constraints struct {
						Text       string `xml:",chardata"`
						Constraint struct {
							Text                   string `xml:",chardata"`
							Field                  string `xml:"field,attr"`
							Scope                  string `xml:"scope,attr"`
							Value                  string `xml:"value,attr"`
							PercentValue           string `xml:"percentValue,attr"`
							Shared                 string `xml:"shared,attr"`
							IncludeChildSelections string `xml:"includeChildSelections,attr"`
							IncludeChildForces     string `xml:"includeChildForces,attr"`
							ID                     string `xml:"id,attr"`
							Type                   string `xml:"type,attr"`
						} `xml:"constraint"`
					} `xml:"constraints"`
				} `xml:"entryLink"`
			} `xml:"entryLinks"`
			SelectionEntryGroups struct {
				Text                string `xml:",chardata"`
				SelectionEntryGroup struct {
					Text       string `xml:",chardata"`
					ID         string `xml:"id,attr"`
					Name       string `xml:"name,attr"`
					Hidden     string `xml:"hidden,attr"`
					Collective string `xml:"collective,attr"`
					Import     string `xml:"import,attr"`
					Modifiers  struct {
						Text     string `xml:",chardata"`
						Modifier struct {
							Text       string `xml:",chardata"`
							Type       string `xml:"type,attr"`
							Field      string `xml:"field,attr"`
							Value      string `xml:"value,attr"`
							Conditions struct {
								Text      string `xml:",chardata"`
								Condition struct {
									Text                   string `xml:",chardata"`
									Field                  string `xml:"field,attr"`
									Scope                  string `xml:"scope,attr"`
									Value                  string `xml:"value,attr"`
									PercentValue           string `xml:"percentValue,attr"`
									Shared                 string `xml:"shared,attr"`
									IncludeChildSelections string `xml:"includeChildSelections,attr"`
									IncludeChildForces     string `xml:"includeChildForces,attr"`
									ChildId                string `xml:"childId,attr"`
									Type                   string `xml:"type,attr"`
								} `xml:"condition"`
							} `xml:"conditions"`
						} `xml:"modifier"`
					} `xml:"modifiers"`
					Constraints struct {
						Text       string `xml:",chardata"`
						Constraint struct {
							Text                   string `xml:",chardata"`
							Field                  string `xml:"field,attr"`
							Scope                  string `xml:"scope,attr"`
							Value                  string `xml:"value,attr"`
							PercentValue           string `xml:"percentValue,attr"`
							Shared                 string `xml:"shared,attr"`
							IncludeChildSelections string `xml:"includeChildSelections,attr"`
							IncludeChildForces     string `xml:"includeChildForces,attr"`
							ID                     string `xml:"id,attr"`
							Type                   string `xml:"type,attr"`
						} `xml:"constraint"`
					} `xml:"constraints"`
					EntryLinks struct {
						Text      string `xml:",chardata"`
						EntryLink []struct {
							Text       string `xml:",chardata"`
							ID         string `xml:"id,attr"`
							Name       string `xml:"name,attr"`
							Hidden     string `xml:"hidden,attr"`
							Collective string `xml:"collective,attr"`
							Import     string `xml:"import,attr"`
							TargetId   string `xml:"targetId,attr"`
							Type       string `xml:"type,attr"`
							Modifiers  struct {
								Text     string `xml:",chardata"`
								Modifier struct {
									Text       string `xml:",chardata"`
									Type       string `xml:"type,attr"`
									Field      string `xml:"field,attr"`
									Value      string `xml:"value,attr"`
									Conditions struct {
										Text      string `xml:",chardata"`
										Condition struct {
											Text                   string `xml:",chardata"`
											Field                  string `xml:"field,attr"`
											Scope                  string `xml:"scope,attr"`
											Value                  string `xml:"value,attr"`
											PercentValue           string `xml:"percentValue,attr"`
											Shared                 string `xml:"shared,attr"`
											IncludeChildSelections string `xml:"includeChildSelections,attr"`
											IncludeChildForces     string `xml:"includeChildForces,attr"`
											ChildId                string `xml:"childId,attr"`
											Type                   string `xml:"type,attr"`
										} `xml:"condition"`
									} `xml:"conditions"`
									ConditionGroups struct {
										Text           string `xml:",chardata"`
										ConditionGroup struct {
											Text       string `xml:",chardata"`
											Type       string `xml:"type,attr"`
											Conditions struct {
												Text      string `xml:",chardata"`
												Condition []struct {
													Text                   string `xml:",chardata"`
													Field                  string `xml:"field,attr"`
													Scope                  string `xml:"scope,attr"`
													Value                  string `xml:"value,attr"`
													PercentValue           string `xml:"percentValue,attr"`
													Shared                 string `xml:"shared,attr"`
													IncludeChildSelections string `xml:"includeChildSelections,attr"`
													IncludeChildForces     string `xml:"includeChildForces,attr"`
													ChildId                string `xml:"childId,attr"`
													Type                   string `xml:"type,attr"`
												} `xml:"condition"`
											} `xml:"conditions"`
										} `xml:"conditionGroup"`
									} `xml:"conditionGroups"`
								} `xml:"modifier"`
							} `xml:"modifiers"`
							Constraints struct {
								Text       string `xml:",chardata"`
								Constraint struct {
									Text                   string `xml:",chardata"`
									Field                  string `xml:"field,attr"`
									Scope                  string `xml:"scope,attr"`
									Value                  string `xml:"value,attr"`
									PercentValue           string `xml:"percentValue,attr"`
									Shared                 string `xml:"shared,attr"`
									IncludeChildSelections string `xml:"includeChildSelections,attr"`
									IncludeChildForces     string `xml:"includeChildForces,attr"`
									ID                     string `xml:"id,attr"`
									Type                   string `xml:"type,attr"`
								} `xml:"constraint"`
							} `xml:"constraints"`
						} `xml:"entryLink"`
					} `xml:"entryLinks"`
				} `xml:"selectionEntryGroup"`
			} `xml:"selectionEntryGroups"`
			CategoryLinks struct {
				Text         string `xml:",chardata"`
				CategoryLink struct {
					Text     string `xml:",chardata"`
					ID       string `xml:"id,attr"`
					Name     string `xml:"name,attr"`
					Hidden   string `xml:"hidden,attr"`
					TargetId string `xml:"targetId,attr"`
					Primary  string `xml:"primary,attr"`
				} `xml:"categoryLink"`
			} `xml:"categoryLinks"`
		} `xml:"selectionEntryGroup"`
	} `xml:"sharedSelectionEntryGroups"`
	SharedRules struct {
		Text string `xml:",chardata"`
		Rule struct {
			Text          string `xml:",chardata"`
			ID            string `xml:"id,attr"`
			Name          string `xml:"name,attr"`
			PublicationId string `xml:"publicationId,attr"`
			Page          string `xml:"page,attr"`
			Hidden        string `xml:"hidden,attr"`
			Description   string `xml:"description"`
		} `xml:"rule"`
	} `xml:"sharedRules"`
	SharedProfiles struct {
		Text    string `xml:",chardata"`
		Profile []struct {
			Text            string `xml:",chardata"`
			ID              string `xml:"id,attr"`
			Name            string `xml:"name,attr"`
			Hidden          string `xml:"hidden,attr"`
			TypeId          string `xml:"typeId,attr"`
			TypeName        string `xml:"typeName,attr"`
			Characteristics struct {
				Text           string `xml:",chardata"`
				Characteristic []struct {
					Text   string `xml:",chardata"`
					Name   string `xml:"name,attr"`
					TypeId string `xml:"typeId,attr"`
				} `xml:"characteristic"`
			} `xml:"characteristics"`
		} `xml:"profile"`
	} `xml:"sharedProfiles"`
}

// GetData fetches the Battlescribe data for the BSData/wh40k repo
func GetData(repo, tag string) ([]*Catalogue, error) {
	log.Infof("getting %s catalogues", repo)

	// clean up
	cleanUp(repo)

	// clone the repo
	err := clone(repo, tag)
	checkErr(err)

	// get the cat files
	files, err := getCatFiles(repo)
	checkErr(err)

	var catalogues []*Catalogue

	// iterate through the cat files and parse them into catalogues
	for _, file := range files {
		log.Infof("Inspecting file %s", file.Name())
		b, err := ioutil.ReadFile(fmt.Sprintf("%s/%s/%s", directory, repo, file.Name()))
		checkErr(err)

		var cat Catalogue
		err = xml.Unmarshal(b, &cat)
		checkErr(err)

		log.Infof("Appending %s", cat.Name)

		catalogues = append(catalogues, &cat)
	}

	// clean up
	cleanUp(repo)

	return catalogues, nil
}

func cleanUp(repo string) {
	if _, err := os.Stat(directory + "/" + repo); !os.IsNotExist(err) {
		err = os.RemoveAll(directory + "/" + repo)
		checkErr(err)
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func clone(repo, tag string) error {
	log.Infof("cloning repo %s/%s", baseDataRepoURL, repo)
	m := sideband.NewMuxer(sideband.Sideband, os.Stdout)
	r, err := git.PlainClone(directory+"/"+repo, false, &git.CloneOptions{
		URL:      baseDataRepoURL + "/" + repo,
		Progress: m,
		Depth:    1,
	})
	if err != nil {
		return err
	}

	ref, err := r.Head()
	if err != nil {
		return err
	}

	log.Infof("checked out at hash (%s): %s", ref.Name(), ref.Hash().String())

	if tag != "" {
		// get the tag's ref
		ref, err = r.Tag(tag)
		if err != nil {
			return err
		}

		// get the worktree
		w, err := r.Worktree()
		if err != nil {
			return err
		}

		// checkout the tag
		err = w.Checkout(&git.CheckoutOptions{
			Hash: ref.Hash(),
		})
		if err != nil {
			return err
		}

		log.Infof("checked out at hash (%s): %s", ref.Name(), ref.Hash().String())
	}

	return nil
}

func getCatFiles(repo string) ([]os.FileInfo, error) {
	fs, err := ioutil.ReadDir(directory + "/" + repo)
	if err != nil {
		return nil, err
	}

	var files []os.FileInfo
	for _, file := range fs {
		isCat := strings.Contains(file.Name(), ".cat")
		log.Infof("Is %s a .cat file? %v", file.Name(), isCat)
		if isCat {
			files = append(files, file)
		}
	}

	return files, nil
}
