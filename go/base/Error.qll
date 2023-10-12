import base.Target
import DataFlow
import base.functions

class Error extends Target::Error::ErrorDef{


    string getHashInfo(){
        result = this.getLocationInfo()
    }
    
}

class ErrorCheck extends Target::Error::ErrorCheck{
    predicate dominatesReturn(Return r){
        Core::checkDomination(this.getTrueGuardNode(), r)
    }

    Error getAssociatedError(){
        result = checkedError
    }

    string getHashInfo(){
       result = this.getLocationInfo()
    }
}

class Return extends Target::Error::CoreReturn{
    FlowNode getAReturnValue(){
        exists(FlowNode f | f.correspondsToExpr(getAReturnedExpr())| result = f)
    }

    string getHashInfo(){
        result = this.getLocationInfo()
    }
}

