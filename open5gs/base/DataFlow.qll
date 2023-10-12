// import base.Target
// import base.Core
// import base.Error
// import base.AKA

// class FlowNode extends Target::DataFlow::DataFlowNode {

//     string getHashInfo(){
//         result = this.getLocationInfo()
//     }
// }

// class CrashRecoveryDataFlow extends Target::DataFlow::CoreTaintConfig{

//     // predicate hasFlow(FlowNode src, FlowNode snk){
//     //     this.(Target::DataFlow::CoreTaintConfig).hasFlow(src, snk)
//     // }


//     override FlowNode defineSource(){
//         result = any(FlowNode d, Core::AkaHandler a | d.correspondsToParamOf(a) | d)
//     }

//     override FlowNode defineSink(){
//         result = any(FlowNode d, Core::AkaFunction a | d.correspondsToParamOf(a) and a.mayLeadToCrash() | d)
//     }

//     override FlowNode defineSanitizer(){
//         result = any(FlowNode d, Core::AkaFunction a | d.correspondsToParamOf(a) and a.recoversFromCrash() | d)
//     }
// }

// class AuthVarDiscovery extends Target::DataFlow::CoreTaintConfig {
//     override FlowNode defineSource(){
//         result = any(FlowNode d, Core::AkaHandler a | d.correspondsToParamOf(a) | d)
//     }

//     override FlowNode defineSink(){
//         result = any(Core::AkaVariable a)
//     }

//     override FlowNode defineSanitizer(){
//         none()
//     }
// }

// class AuthVarProtected extends Target::DataFlow::CoreDataFlowConfig2{
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
//     Target::DataFlow::hasLocalTaint(src, sink)
// }

// predicate hasLocalFlow(FlowNode src, FlowNode sink){
//     Target::DataFlow::hasLocalFlow(src, sink)
// }
