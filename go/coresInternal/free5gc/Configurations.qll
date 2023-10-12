private import BaseLanguage::BaseLanguageLib as BL
import Nodes
import ErrorDef
class CoreDataFlowConfig extends BL::DataFlow::Configuration {
    CoreDataFlowConfig() { this = "CoreDataFlowConfig" }
    
    override predicate isSource(BL::DataFlow::Node node) {
        node = defineSource()
    }
    
    override predicate isSink(BL::DataFlow::Node node) {
        node = defineSink()
    }

    override predicate isBarrier(BL::DataFlow::Node node){
        node = defineBarrier()
    }

    abstract DataFlowNode defineSource();
    abstract DataFlowNode defineSink();
    abstract DataFlowNode defineBarrier();

    override predicate isAdditionalFlowStep(BL::DataFlow::Node fromNode, BL::DataFlow::Node toNode) {
        fromNode = toNode.(CoreCallNode).getReceiver() and fromNode.getType() instanceof BL::ErrorType
        and toNode.toString() = "call to Error"
    }


}


class CoreTaintConfig extends BL::DataFlow::TaintTracking::Configuration {
    CoreTaintConfig() { this = "CoreTaintConfig" }
    
    override predicate isSource(BL::DataFlow::Node node) {
        node = defineSource()
    }
    
    override predicate isSink(BL::DataFlow::Node node) {
        node = defineSink()
    }

    override predicate isSanitizer(BL::DataFlow::Node node){
        node = defineSanitizer()
    }

    override predicate isAdditionalTaintStep(BL::DataFlow::Node fromNode, BL::DataFlow::Node toNode) {
        any(BL::DataFlow::Write w).writesComponent(toNode.(BL::DataFlow::PostUpdateNode).getPreUpdateNode(), fromNode)
        or
        (
            fromNode.asExpr().getParent() = toNode.asExpr() 
            and exists(BL::DataFlow::Write w | w.getRhs().asExpr() = toNode.asExpr())
            and fromNode.asExpr().toString() = "index expression"
        )    
        or
        (
            toNode instanceof BL::DataFlow::BinaryOperationNode 
            and toNode.(BL::DataFlow::BinaryOperationNode).getRightOperand() = fromNode
            and fromNode.asExpr().toString() = "index expression"
        )
        // or
        // (
        //     exists(BL::DataFlow::CallNode c, BL::StructLit s| fromNode = c.getReceiver()
        //     and c.asExpr().getParent() = s.getAChildExpr() and toNode = BL::DataFlow::exprNode(s))
        // )
    }

    abstract DataFlowNode defineSource();
    abstract DataFlowNode defineSink();
    abstract DataFlowNode defineSanitizer();
    // abstract predicate defineAdditionalTaintStep(DataFlowNode fromNode, DataFlowNode toNode);

    
}

class CoreTaintConfig2 extends BL::DataFlow2::TaintTracking::Configuration {
    CoreTaintConfig2() { this = "CoreTaintConfig2" }
    
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

class CoreDataFlowConfig2 extends BL::DataFlow2::Configuration {
    CoreDataFlowConfig2() { this = "CoreDataFlowConfig2" }
    
    override predicate isSource(BL::DataFlow::Node node) {
        node = defineSource()
    }
    
    override predicate isSink(BL::DataFlow::Node node) {
        node = defineSink()
    }

    override predicate isBarrier(BL::DataFlow::Node node){
        node = defineBarrier()
    }

    abstract DataFlowNode defineSource();
    abstract DataFlowNode defineSink();
    abstract DataFlowNode defineBarrier();


}


class StructAndSliceTaint extends BL::DataFlow::TaintTracking::AdditionalTaintStep{
    
    override predicate step(BL::DataFlow::Node fromNode, BL::DataFlow::Node toNode){
        any(BL::DataFlow::Write w).writesComponent(toNode.(BL::DataFlow::PostUpdateNode).getPreUpdateNode(), fromNode)
        or
        (
            fromNode.asExpr().getParent() = toNode.asExpr() 
            and exists(BL::DataFlow::Write w | w.getRhs().asExpr() = toNode.asExpr())
            and fromNode.asExpr().toString() = "index expression"
        )    
        or
        (
            toNode instanceof BL::DataFlow::BinaryOperationNode 
            and toNode.(BL::DataFlow::BinaryOperationNode).getRightOperand() = fromNode
            and fromNode.asExpr().toString() = "index expression"
        )
        or
        (
            exists(BL::DataFlow::CallNode c, BL::StructLit s| fromNode = c
            and c.asExpr().getParent() = s.getAChildExpr() and toNode = BL::DataFlow::exprNode(s))
        )
        or
        (fromNode = toNode.(CoreCallNode).getReceiver() and fromNode.getType() instanceof BL::ErrorType)
    }
}