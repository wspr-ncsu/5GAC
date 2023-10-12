// import base.Target
// import DataFlow2
// import base.functions

// class Error extends Target::Error2::ErrorDef{


//     string getHashInfo(){
//         result = this.getLocationInfo()
//     }
// }

// class ErrorCheck extends Target::Error2::ErrorCheck{
//     predicate dominatesReturn(Return r){
//         Core::checkDomination(this.getTrueGuardNode(), r)
//     }

//     Error getAssociatedError(){
//         result = checkedError
//     }

//     string getHashInfo(){
//        result = this.getLocationInfo()
//     }
// }

// class Return extends Target::Error2::CoreReturn{
//     FlowNode getAReturnValue(){
//         exists(FlowNode f | f.correspondsToExpr(getAReturnedExpr()) | result = f)
//     }

//     string getHashInfo(){
//         result = this.getLocationInfo()
//     }
// }

