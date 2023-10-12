// import base.Target
// import base.Core


// class Log extends Target::DataFlow::LogCallDef {
// }


// class LogSearch extends Target::DataFlow::CoreDataFlowConfig {

//     override Core::FlowNode defineSource(){
//         result = any(Core::FlowNode a)
//     }

//     override Core::FlowNode defineSink(){
//         result = any(Log l | | l.getAMessageComponent())
//     }

//     override Core::FlowNode defineBarrier(){
//         result = any(Log l, Core::FlowNode f | f.correspondsToParamOf(l.(Core::CallNode).getTarget()) | f)
//         // none()
//     }
// }