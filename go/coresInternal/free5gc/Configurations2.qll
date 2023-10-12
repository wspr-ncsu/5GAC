private import BaseLanguage::BaseLanguageLib as BL
import Nodes2

class CoreDataFlowConfig extends BL::DataFlow2::Configuration {
    CoreDataFlowConfig() { this = "CoreDataFlowConfig" }
    
    override predicate isSource(BL::DataFlow2::Node node) {
        node = defineSource()
    }
    
    override predicate isSink(BL::DataFlow2::Node node) {
        node = defineSink()
    }

    override predicate isBarrier(BL::DataFlow2::Node node){
        node = defineBarrier()
    }

    abstract DataFlowNode defineSource();
    abstract DataFlowNode defineSink();
    abstract DataFlowNode defineBarrier();


}


class CoreTaintConfig extends BL::DataFlow2::TaintTracking::Configuration {
    CoreTaintConfig() { this = "CoreTaintConfig" }
    
    override predicate isSource(BL::DataFlow2::Node node) {
        node = defineSource()
    }
    
    override predicate isSink(BL::DataFlow2::Node node) {
        node = defineSink()
    }

    override predicate isSanitizer(BL::DataFlow2::Node node){
        node = defineSanitizer()
    }

    override predicate isAdditionalTaintStep(BL::DataFlow2::Node fromNode, BL::DataFlow2::Node toNode) {
        any(BL::DataFlow2::Write w).writesComponent(toNode.(BL::DataFlow2::PostUpdateNode).getPreUpdateNode(), fromNode)
        or
        (
            fromNode.asExpr().getParent() = toNode.asExpr() 
            and exists(BL::DataFlow2::Write w | w.getRhs().asExpr() = toNode.asExpr())
            and fromNode.asExpr().toString() = "index expression"
        )    
        or
        (
            toNode instanceof BL::DataFlow2::BinaryOperationNode 
            and toNode.(BL::DataFlow2::BinaryOperationNode).getRightOperand() = fromNode
            and fromNode.asExpr().toString() = "index expression"
        )
    }

    abstract DataFlowNode defineSource();
    abstract DataFlowNode defineSink();
    abstract DataFlowNode defineSanitizer();
    // abstract predicate defineAdditionalTaintStep(DataFlowNode fromNode, DataFlowNode toNode);

    
}