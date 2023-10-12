import Functions
import Expr
private import BaseLanguage::BaseLanguageLib as BL

class CoreHandlers extends CoreCall {

    int path_index;

    CoreHandlers(){
        (this.getTargetCore().getParameter(_).getType().getName() = "ogs_sbi_stream_t *"
        or
        this.getTargetCore().getParameter(_).getType().getName() = "ogs_sbi_message_t *")
    and not this.getTargetCore().getName() = "ogs_sbi_server_send_error"
    and (isHandlerRootSwitchDescendant(this.getEnclosingBlock()) 
        or (not isHandlerRootSwitchDescendant(this.getEnclosingBlock()) 
                and isEventDescendant(this.getEnclosingBlock())))
    and
    path_index = max(int i | i = 0 or i = any(HandlerResourceAccess h 
        | h.getEnclosingBlock() = (this.getEnclosingBlock()).getEnclosingBlock*() 
        | h.getArrayIndexInt()))

    }

    
    string getHttpTypeInternal() {
        result = getHttpType(this.getEnclosingBlock())
    }


    string getSerivceName(){
        exists(BL::IfStmt i,  BL::StringLiteral e, HandlerRootSwitch h 
            | isBlockDescendent(i.getThen(), this.getEnclosingBlock())
                and e.getEnclosingStmt() = i and h.getService() = e
            | result = e.toString())
    }

    CoreFunction asFunction(){
        result = this.getTargetCore()
    }

    // /**
    //  * Core-specific logic to determine what URL reaches this handler.
    //  * Returns: String of handler URL
    //  */
    string getPathInternal(){
        // if(count(string s | s = getRoutingString(this.getEnclosingBlock()).splitAt("/") | s) - 1) < path_index then
        //     result = concat(string str, int x 
        //     | x in [0 .. (path_index - count(string s | s = getRoutingString(this.getEnclosingBlock()).splitAt("/")))] 
        //         and str = "unspecified/" 
        //     | str) + getRoutingString(this.getEnclosingBlock())
        // else
        //         result = getRoutingString(this.getEnclosingBlock())
    result = concat(int x, string s | x in [0 .. path_index] and s = getPathAtIndex(x, this.getEnclosingBlock()) | s , "/" order by x)

    }

    string getNFInternal() {result = "placeholder"}

    string getName(){
        result =this.getTargetCore().getName()
    }


    
}