package restCallsToCAEngine

import (
	"github.com/jlambert68/FenixTestInstructionsDataAdmin/SubCustody/FangEngineClassesAndMethods"
	"github.com/jlambert68/FenixTestInstructionsDataAdmin/SubCustody/TestInstructions"
	"github.com/jlambert68/FenixTestInstructionsDataAdmin/TypeAndStructs"
)

// FangEngineRestApiMessageStruct
// Used when converting a TestInstruction into data to be used in RestCall to FangEngine
type FangEngineRestApiMessageStruct struct {
	FangEngineClassNameNAME           FangEngineClassesAndMethods.FangEngine_ClassName_Name_SC_Type
	FangEngineMethodNameNAME          FangEngineClassesAndMethods.FangEngine_MethodName_Name_SC_Type
	FangEngineExpectedToBePassedValue TypeAndStructs.AttributeValueAsStringType
	TestInstructionAttribute          []TestInstructionAttributesUuidAndValueStruct
	FangAttributes                    map[TypeAndStructs.TestInstructionAttributeUUIDType]*FangEngineClassesAndMethods.FangEngineAttributesStruct

	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// TestInstructionAttributesUuidAndValueStruct
// Holds one Attribute UUID, Name and Value
type TestInstructionAttributesUuidAndValueStruct struct {
	TestInstructionAttributeUUID          TypeAndStructs.TestInstructionAttributeUUIDType
	TestInstructionAttributeName          TypeAndStructs.TestInstructionAttributeNameType
	TestInstructionAttributeValueAsString TypeAndStructs.AttributeValueAsStringType
}

var (
	// All TestInstruction-data for 'Sub Custody'
	allTestInstructions_SC TestInstructions.AllTestInstructions_SC_TestCaseSetUpStruct

	// All Attributes-data for 'Sub Custody' as map
	testInstructionAttributesMap map[TypeAndStructs.TestInstructionAttributeUUIDType]*TypeAndStructs.TestInstructionAttributeStruct

	// All FangEngineData for 'Sub Custody' as map
	fangEngineClassesMethodsAttributesMap map[TypeAndStructs.OriginalElementUUIDType]*FangEngineClassesAndMethods.FangEngineClassesMethodsAttributesStruct
)
