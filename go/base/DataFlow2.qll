// import base.Target
// import base.Core
// import base.Error2
// import base.AKA

// class FlowNode extends Target::DataFlow2::DataFlowNode {

//     string getHashInfo(){
//         result = this.getLocationInfo()
//     }
// }


// class AuthVarProtected extends Target::DataFlow2::CoreDataFlowConfig{
//     override FlowNode defineSource(){
//         result = any(AuthVariable a)
//     }

//     override FlowNode defineSink(){
//         exists(FlowNode n, AuthVariable a, ErrorCheck e| n.readsVar(a) 
//             and n.asInstruction() = e.getTrueGuardNode().getASuccessor*()| result = n)
//     }

//     override FlowNode defineBarrier(){
//         exists(FlowNode n, AuthVariable a, ErrorCheck e| n.writesVar(a) | result = n)
//     }

// }


// predicate hasLocalTaint(FlowNode src, FlowNode sink){
//     Target::DataFlow2::hasLocalTaint(src, sink)
// }

// predicate hasLocalFlow(FlowNode src, FlowNode sink){
//     Target::DataFlow2::hasLocalFlow(src, sink)
// }
