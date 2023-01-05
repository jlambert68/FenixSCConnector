package restCallsToCAEngine

import (
	"github.com/jlambert68/FenixTestInstructionsDataAdmin/SubCustody/FangEngineClassesAndMethods"
	"github.com/jlambert68/FenixTestInstructionsDataAdmin/SubCustody/TestInstructions"
	"github.com/jlambert68/FenixTestInstructionsDataAdmin/TypeAndStructs"
)

// InitiateRestCallsToCAEngine
// Do all initiation to have restEngine be able to do RestCalls to Sub Custodys FangEngine
func InitiateRestCallsToCAEngine() {

	// Load all TestInstruction-data for 'Sub Custody'
	allTestInstructions_SC = TestInstructions.InitiateAllTestInstructionsForSC()

	// Initiate map-objects
	testInstructionAttributesMap = make(map[TypeAndStructs.TestInstructionAttributeUUIDType]*TypeAndStructs.TestInstructionAttributeStruct)
	fangEngineClassesMethodsAttributesMap = make(map[TypeAndStructs.OriginalElementUUIDType]*FangEngineClassesAndMethods.FangEngineClassesMethodsAttributesStruct)

	// Convert TestInstruction-data for 'Sub Custody' into map-objects
	convertTestInstructionDataIntoMapStructures()

}

// Convert TestInstruction-data for 'Sub Custody' into map-objects
func convertTestInstructionDataIntoMapStructures() {

	// Loop TestInstructionsAttributes and create Map
	for _, testInstructionsAttribute := range allTestInstructions_SC.TestInstructionAttribute {
		var tempTestInstructionsAttribute TypeAndStructs.TestInstructionAttributeStruct

		tempTestInstructionsAttribute = testInstructionsAttribute
		testInstructionAttributesMap[tempTestInstructionsAttribute.TestInstructionAttributeUUID] = &tempTestInstructionsAttribute
	}

	// Loop FangEngineClassesMethodsAttributes and create Map
	for _, fangEngineClassesMethodsAttribute := range allTestInstructions_SC.FangEngineClassesMethodsAttributes {
		var tempFangEngineClassesMethodsAttribute FangEngineClassesAndMethods.FangEngineClassesMethodsAttributesStruct

		tempFangEngineClassesMethodsAttribute = fangEngineClassesMethodsAttribute
		fangEngineClassesMethodsAttributesMap[tempFangEngineClassesMethodsAttribute.TestInstructionOriginalUUID] = &tempFangEngineClassesMethodsAttribute
	}

}
