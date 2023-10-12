import base.Target
import coreapi.Handler
import base.Core
import coreapi.CoreAPI
import utils.Utils

class ServiceCall extends Target::CoreServiceCall {

    string target_nf;
    string path;
    Handler target_handler;
    string http_type;
    string path_root;
    Target::CoreCall call;
    string file;

    ServiceCall() {
        this instanceof Target::CoreServiceCall 
        and
        path = getPathCore() and 
        target_nf = getTargetNFCore() and 
        target_handler = getTargetHandlerCore() and
        http_type = getHttpTypeCore() and
        path_root = getPathRootCore() and
        this.isCalled() and
        call = this.getCallCore() and
        file = call.getFile().getRelativePath()
    }


    string getPath() {
        result = path
    }

    string getTargetNF() {
        result = target_nf
    }

    Handler getTargetHandler() {
        result = target_handler
    }

    string getHttpType() {
        result = http_type
    }

    string getPathRoot(){
        result = path_root
    }
    
    string getColonizedPath(){
        result = this.getColonizedPathCore()
    }

    // string getSourceHandlerName() {
    //     result = source_handler.getName()
    // }

    // Handler getSourceHandler() {
    //     result = source_handler
    // }

    Core::Call getCall(){
        result = call
    }

    string getSourceNF(){
        result = getSourceNFCore()
    }


    

}