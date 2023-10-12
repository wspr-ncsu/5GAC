import Functions
import Handlers
private import BaseLanguage::BaseLanguageLib as BL

class SBIResourceAssignment extends BL::Assignment {
    BL::ArrayExpr resource_array_access;

    SBIResourceAssignment(){
        resource_array_access = this.getLValue().(BL::ArrayExpr)
        and
        resource_array_access.getArrayBase().toString() = "component"
        and
        resource_array_access.getArrayBase().(BL::ValueFieldAccess)
            .getQualifier().toString() = "resource"
        and
        resource_array_access.getArrayBase().(BL::ValueFieldAccess)
            .getQualifier().(BL::ValueFieldAccess).getQualifier()
            .getType().getName() ="ogs_sbi_header_t"
    }

    int getArrayIndex(){
        result = resource_array_access.getArrayOffset().toString().toInt()
    }

    string getResourceValue(){
        if this.getRValue() instanceof BL::StringLiteral then (
            result = this.getRValue().getValue()
        )
        else if this.getRValue() instanceof BL::PointerFieldAccess then (
            result = "$Pointer:"+this.getRValue().(BL::PointerFieldAccess).getTarget().toString()
        )
        else if this.getRValue() instanceof BL::FunctionCall then (
            result = "$Call:"+this.getRValue().(BL::FunctionCall).getArgument(_).(BL::PointerFieldAccess).getTarget().toString()
            or
            result = "$Call:"+this.getRValue().(BL::FunctionCall).getArgument(_).(BL::AddressOfExpr).getAddressable().toString()
        )
        else if this.getRValue() instanceof BL::ArrayExpr then (
            this.getRValue().(BL::ArrayExpr).getArrayBase().toString() = "component"
            and
            this.getRValue().(BL::ArrayExpr).getArrayBase().(BL::ValueFieldAccess)
                .getQualifier().(BL::ValueFieldAccess).getQualifier()
                .getType().getName() ="ogs_sbi_header_t"
            and
            exists(SBIResourceAssignment s | s.getEnclosingFunction() = this.getEnclosingFunction() 
                and s.getArrayIndex() = this.getRValue().(BL::ArrayExpr).getArrayOffset().toString().toInt() |
                result = s.getResourceValue())
        )
        else
            result = this.getRValue().toString()
    }

    BL::ArrayExpr getArrayExpr(){
        result = resource_array_access
    }

}


class SBIHttpAssignment extends BL::Assignment{
    SBIHttpAssignment(){
        this.getLValue().(BL::ValueFieldAccess).toString() = "method"
        and
        this.getLValue().(BL::ValueFieldAccess).getQualifier().toString() = "h"
        and
        this.getLValue().(BL::ValueFieldAccess).getQualifier().getType().getName() = "ogs_sbi_header_t"
    }

    string getMethod(){
        result = this.getRValue().getValue()
    }
}


class SBIServiceAssignment extends BL::Assignment {
    SBIServiceAssignment(){
        this.getLValue().(BL::ValueFieldAccess).toString() = "name"
        and
        this.getLValue().(BL::ValueFieldAccess).getQualifier().toString() = "service"
        and
        this.getLValue().(BL::ValueFieldAccess).getQualifier().(BL::ValueFieldAccess).getQualifier().getType().getName() = "ogs_sbi_header_t"
    }

    string getService(){
        result = this.getRValue().getValue()
    }

    string getNFTarget(){
        result = getService().regexpCapture("n([\\w]+)-[\\w]+", 1)
    }
}

class SBIVerAssignment extends BL::Assignment {
    SBIVerAssignment(){
        this.getLValue().(BL::ValueFieldAccess).toString() = "version"
        and
        this.getLValue().(BL::ValueFieldAccess).getQualifier().toString() = "api"
        and
        this.getLValue().(BL::ValueFieldAccess).getQualifier().(BL::ValueFieldAccess).getQualifier().getType().getName() = "ogs_sbi_header_t"
    }

    string getApiVersion(){
        result = this.getRValue().getValue()
    }

}


class HandlerRootSwitch extends BL::ConditionalExpr {
    HandlerRootSwitch(){
        this.getAnOperand().toString() = "name"
    and this.getAnOperand().(BL::ValueFieldAccess).getQualifier().toString() = "service"
    and this.getAnOperand().(BL::ValueFieldAccess).getQualifier().(BL::ValueFieldAccess).getQualifier().getType().getName() = "ogs_sbi_header_t"
    and this.getEnclosingBlock().getParent().getBasicBlock().getStart() instanceof BL::SwitchStmt
    }

    BL::BlockStmt getSwitchBlock(){
        result = this.getEnclosingBlock()
    }

    BL::StringLiteral getService(){
        exists(BL::StringLiteral s | s.getEnclosingBlock().getEnclosingBlock() = getSwitchBlock() and result = s)
    }
}

class SBIResourceArrayExpr extends BL::ArrayExpr {
    SBIResourceArrayExpr(){
        this.getArrayBase().toString() = "component"
        and
        this.getArrayBase().(BL::ValueFieldAccess).getQualifier().toString() = "resource"
        and
        this.getArrayBase().(BL::ValueFieldAccess).getQualifier().(BL::ValueFieldAccess).getQualifier().getType().getName() = "ogs_sbi_header_t"
    }
}

class SBIHttpTypeExpr extends BL::ValueFieldAccess{
    SBIHttpTypeExpr(){
        this.toString() = "method"
        and
        this.getQualifier().getType().getName() = "ogs_sbi_header_t"
    }
}

class HandlerResourceAccess extends SBIResourceArrayExpr {
    HandlerResourceAccess(){
        isHandlerRootSwitchDescendant(this.getEnclosingBlock())
    }

    string getComparedValue(){
        exists(BL::StringLiteral s | s.getEnclosingBlock().getEnclosingBlock() = this.getEnclosingBlock() | result = s.getValue())
    }

    int getArrayIndexInt(){
        result = this.getArrayOffset().toString().toInt()
    }

}

class HandlerHttpTypeAccess extends SBIHttpTypeExpr{
    HandlerHttpTypeAccess(){
        isHandlerRootSwitchDescendant(this.getEnclosingBlock())
    }

    string getComparedValue(){
        exists(BL::StringLiteral s | s.getEnclosingBlock().getEnclosingBlock() = this.getEnclosingBlock() | result = s.getValue())
    }

}

// class HandlerHttpTypeAccess extends 


predicate isHandlerRootSwitchDescendant(BL::BlockStmt b){
    exists(HandlerRootSwitch h | h.getSwitchBlock() = b.getEnclosingBlock*())
}

string getRoutingString(BL::BlockStmt b){
    result = concat(string s 
        | exists(BL::IfStmt i,  BL::StringLiteral e, HandlerResourceAccess h 
            | isBlockDescendent(i.getThen(), b.getEnclosingBlock*())
        and e.getEnclosingStmt() = i and h.getComparedValue() = e.toString()| s = e.toString())
        | s, "/")
}



bindingset [x]
string getPathAtIndex(int x, BL::BlockStmt b){
    exists(BL::IfStmt i, HandlerResourceAccess h, BL::StringLiteral e
        | isBlockDescendent(i.getThen(), b.getEnclosingBlock*())
        and e.getEnclosingStmt() = i and h.getComparedValue() = e.toString()
        and h.getArrayIndexInt() = x 
        and h.isInMacroExpansion() | result = e.getValue())
    or
    (not exists(BL::IfStmt i, HandlerResourceAccess h, BL::StringLiteral e
        | isBlockDescendent(i.getThen(), b.getEnclosingBlock*())
        and e.getEnclosingStmt() = i and h.getComparedValue() = e.toString()
        and h.getArrayIndexInt() = x
        and h.isInMacroExpansion())  
    and result = "UNSPECIFIED")
}


// exists(BL::IfStmt i,  BL::StringLiteral e, HandlerRootSwitch h | isBlockDescendent(i.getThen(), this.getEnclosingBlock())
// and e.getEnclosingStmt() = i and h.getService() = e| result = e.toString())

string getHttpType(BL::BlockStmt b){
    // if 
    not exists(HandlerHttpTypeAccess h | h.getEnclosingBlock() = b.getEnclosingBlock*()) //then
        and result = "Unspecified"
    // else
    or
    exists(BL::IfStmt i,  BL::StringLiteral e, HandlerHttpTypeAccess h 
        | isBlockDescendent(i.getThen(), b.getEnclosingBlock*())
    and e.getEnclosingStmt() = i and h.getComparedValue() = e.toString()| result = e.toString())
}

HandlerRootSwitch getHandlerRootSwitch(BL::BlockStmt b){
    exists(HandlerRootSwitch h | h.getSwitchBlock() = b.getEnclosingBlock*() | result = h)

}

predicate isBlockDescendent(BL::BlockStmt parent, BL::BlockStmt child){
    parent = child.getEnclosingBlock*()
} 

predicate isStmtDescendent(BL::Stmt parent, BL::Stmt child) {
    parent = child.getEnclosingBlock().getEnclosingStmt()
    
}

predicate isFunctionDescendent(CoreFunction child, CoreFunction parent){
    child.getCallerCore*().getCallerCore().(CoreFunction) = parent
}

predicate isEventDescendant(BL::BlockStmt b){
        exists(BL::SwitchCase c, BL::SwitchStmt s 
            | c.getSwitchStmt() = s
                and s.getControllingExpr().toString() = "id" 
                and (c.getExpr().toString().matches("%SBI_CLIENT") or c.getExpr().toString().matches("%SBI_SERVER"))
            | c.getAStmt() = b.getEnclosingBlock*())
}