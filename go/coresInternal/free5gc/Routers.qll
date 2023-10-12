private import BaseLanguage::BaseLanguageLib as BL
import Functions

class CoreRouter extends BL::StructLit {
    string handler_name;
    CoreRouter() {
        this = any(BL::SliceLit routes | 
            routes.getType().getName() = "Routes" |
            routes.getElement(_))
            and handler_name = this.getElement(3).toString().replaceAll("\"", "")
            and not handler_name = "Index"
    }


    string getHttpTypeCore() {
        (this.getElement(1) instanceof CoreCall and 
        result = this.getElement(1).(CoreCall).getArgument(0).getExactValue())
        or
        (this.getElement(1) instanceof BL::StringLit and 
        result = this.getElement(1).toString().replaceAll("\"", ""))
    }

    string getPathCore() {
        this.getElement(2) instanceof BL::StringLit and 
        result = this.getElement(2).toString().replaceAll("\"", "")
    }

    string getHandlerNameCore(){
        result = handler_name
    }

    string getRouterRootPath(){
        exists(CoreCall c | 
            c.getTargetCore().getName() = "Group" 
            and c.getFile() = this.getFile()
            and c.getEnclosingFunction().getName() = "AddService" |
            result = c.getArgument(0).getExactValue())
    }
}


class CoreRoutingFunction extends BL::IfStmt{
    BL::CallExpr routing_call;
    BL::BlockStmt then_block;



    CoreRoutingFunction(){  
        this = any(BL::IfStmt i, BL::CallExpr c |
            c.getTarget().getName()="Any"
            and i.getParent() = c.getArgument(1).(BL::FunctionName).getTarget().getBody()
            | i) and then_block = this.getThen() and 
            routing_call = any(BL::CallExpr c 
                | c.getTarget().getName() = "Any"
                and c.getArgument(1).(BL::FunctionName).getTarget().getBody() = this.getParent())
    }


    string getHttpTypeCore() {
        exists(CoreCall c
             | 
            this.getCond().getChild(_).(BL::EqlExpr).getAnOperand() = c
            and c.getTargetCore().getName() = "ToUpper"
            | result = c.getArgument(0).getExactValue())
        or
        // For Single length cases
        exists(CoreCall c 
            | this.getCond().(BL::EqlExpr).getAnOperand() = c
            and c.getTargetCore().getName() = "ToUpper"
            | result = c.getArgument(0).getExactValue())

    }

    private string getOp(){
        exists(BL::VariableName c, string v, BL::EqlExpr e
            | this.getCond().getChild(_) = e 
            and e.getAnOperand() = c
            and this.getCond().getChild(_).(BL::EqlExpr).getAnOperand().(BL::StringLit).getExactValue() = v
            | result = v)
    }

    private string getOpKey(){

            this.getCond().getChild(_).(BL::EqlExpr)
                .getAnOperand().(BL::VariableName).getTarget().getAWrite()
                .getRhs().(CoreCallNode).getArgument(0).getExactValue() = result
    }

    private string replacePathAlignment(){
        exists(CoreCall c, BL::KeyValueExpr kv, BL::KeyValueExpr kv2, BL::AssignStmt a,
            string key, string value
            | a.getParent() = routing_call.getArgument(1).(BL::FunctionName).getTarget().getBody() // Assignment is inside this function
            and a.getRhs() = c // Enforce CoreCall
            and c.getTargetCore().getName() = "append" // Check is call to append
            and c.getArgument(1).(BL::StructLit).getElement(0) = kv
            and c.getArgument(1).(BL::StructLit).getElement(1) = kv2 // Enforces StructLit argument and gets key-value expression
            and 
            kv.getKey().toString() = "Key" and kv2.getKey().toString() = "Value" and
                    key = kv.getValue().toString().replaceAll("\"", "")
                    and
                    value = kv2.getValue().(CoreCall).getArgument(0).toString().replaceAll("\"", "")
            | result = getPathBase().replaceAll(value, key+"$")
        )
    }

    private string replaceOpOnly(){
        result = getPathBase().regexpReplaceAll(":"+getOpKey()+"\\b(?!\\$)", getOp()+"\\$")
    }

    private string replaceOpOneAppend(){
        exists(CoreCall c, BL::KeyValueExpr kv, BL::KeyValueExpr kv2, BL::AssignStmt a,
            string key, string value
            | a.getParent() = then_block // Assignment is inside this function
            and a.getRhs() = c // Enforce CoreCall
            and c.getTargetCore().getName() = "append" // Check is call to append
            and c.getArgument(1).(BL::StructLit).getElement(0) = kv
            and c.getArgument(1).(BL::StructLit).getElement(1) = kv2 // Enforces StructLit argument and gets key-value expression
            and 
            kv.getKey().toString() = "Key" and kv2.getKey().toString() = "Value" and
                    key = kv.getValue().toString().replaceAll("\"", "")
                    and
                    value = kv2.getValue().(CoreCall).getArgument(0).toString().replaceAll("\"", "")
            | result = replaceOpOnly().replaceAll(value, key+"$")
        )
    }

    private string replaceOpTwoAppend(){
        exists(CoreCall c, BL::KeyValueExpr kv, BL::KeyValueExpr kv2, BL::AssignStmt a,
            string key2, string value2, string onerep
            | a.getParent() = then_block // Assignment is inside this function
            and a.getRhs() = c // Enforce CoreCall
            and c.getTargetCore().getName() = "append" // Check is call to append
            and c.getArgument(1).(BL::StructLit).getElement(0) = kv
            and c.getArgument(1).(BL::StructLit).getElement(1) = kv2 // Enforces StructLit argument and gets key-value expression
            and 
            kv.getKey().toString() = "Key" and kv2.getKey().toString() = "Value" and
                    key2 = kv.getValue().toString().replaceAll("\"", "")
                    and
                    value2 = kv2.getValue().(CoreCall).getArgument(0).toString().replaceAll("\"", "")
                    and onerep = replaceOpOneAppend() 
                    and onerep.regexpFind(value2+"\\b(?!\\$)", _, _) = value2
            | result = onerep.regexpReplaceAll(value2+"\\b(?!\\$)", key2+"\\$")
        )
    }


    string getPathCore() {
        if count(string s | s = replaceOpTwoAppend().splitAt("$") | s ) > 3 then
            result = replaceOpTwoAppend().replaceAll("$", "").replaceAll("am-data/:thirdLayer", "am-data/sor-ack")
        else if count(string s | s = replaceOpOneAppend().splitAt("$") | s) > 2 then
            result = replaceOpOneAppend().replaceAll("$", "")
        else if count(string s | s = replaceOpOnly().splitAt("$") | s) > 1 then
            result = replaceOpOnly().replaceAll("$", "")
        else if count(string s | s = replacePathAlignment().splitAt("$") | s) > 1 then
            result = replacePathAlignment().replaceAll("$", "")
        else
            result = getPathBase()
    }


    string getPathBase() {
        exists(BL::SsaVariable s, BL::Write w |
            routing_call.getArgument(0).(BL::VariableName).getTarget().(BL::SsaSourceVariable) = s.getSourceVariable() 
            and w.definesSsaVariable(s, _) | result = w.getRhs().getExactValue())
    }
    string getHandlerNameCore(){
        if then_block.getChildStmt(then_block.getNumChildStmt()-2).getChild(0).getParent().getParent()
        = then_block then
            result = then_block.getChildStmt(then_block.getNumChildStmt()-2).getChild(0).(CoreCall).getTargetCore().getName()
        else
            result = "Invalid"
    }

    string getRouterRootPath(){
        exists(BL::CallExpr c | 
            c.getTarget().getName() = "Group" 
            and c.getEnclosingFunction().getBody() = routing_call.getEnclosingFunction().getBody() |
            result = c.getArgument(0).getExactValue())
    }

    BL::CallExpr getRoutingFunctionCall(){
        result = routing_call
    }

    BL::IfStmt getIfStmtCore(){
        result = this
    }

    BL::BlockStmt getThenBlock(){
        result = then_block
    }

}