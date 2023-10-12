import base.Target

class Handler extends Target::CoreHandlers {

    // string nf;
    string url;
    string http_type;
    string name;
    string nf;

    Handler() {
        this instanceof Target::CoreHandlers 
        and nf = this.getNFInternal()
        and url = this.getPathInternal() 
        and http_type = this.getHttpTypeInternal()
        and name = this.getName()
    }

    string getNF() { result = nf }

    string getPath() { result = url }

    string getHttpType() { result = http_type }

}