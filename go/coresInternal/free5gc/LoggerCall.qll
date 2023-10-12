private import BaseLanguage::BaseLanguageLib as BL
import Nodes
/** 
 * LogCallDef is used as the sink for Dataflow to log queries
 * 
 * Core-specific class defining what to consider a call to a log write
 * In Go, we can use LoggerCall
 *  Other languages will be different definitions.
*/
class LogCallDef extends DataFlowNode, BL::LoggerCall {

    DataFlowNode getLoggedValue(){
        result = this.getAMessageComponent()
    }
}


// Extends LoggerCall::Range to include free5gc specific log calls

class AdditionalLoggerCall extends BL::LoggerCall::Range, CoreCallNode {
    AdditionalLoggerCall() {
        (this.getTarget().getName().matches(["Errorf", "Fprintln", "ReportError", "Println", "Fprintf", "Printf", "PrintResult", 
        "Errorln", "log", "Log", "Logf", "errorf", "log", "Log", "Logf", "Fprint", "PrintFn","LogFn","logf", "Logln", "Print",
        "printCriticalityDiagnostics","printAndGetCause", "Debugln", "Infoln", "Println", "Warnln", "Warningln", "Fatalln", 
        "Panicln", "Debug", "Info", "Print", "Warn", "Warning", "Fatal", "Panic", "Debugf", "Infof", "Printf", "Warnf", 
        "Warningf", "Errorf", "Fatalf", "Panicf", "WithField", "WithFields", "WithError"]) 
        
        or this.getTarget().hasQualifiedName("io", ["WriteString", "Writer", "MultiWriter", "WriterAt", "WriterTo", "Copy",
                                                    "CopyBuffer", "CopyN"]))
        and not this.getEnclosingCallable().getName() = "ReportError"

    }

    override BL::DataFlow::Node getAMessageComponent() {
        result = this.getAnArgument()
    }
}
