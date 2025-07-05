export namespace backend {
	
	export class CheckParams {
	    ProxyList: string[];
	    ProxyType: string;
	    Endpoint: string;
	    Threads: number;
	    UpstreamProxy?: string;
	    UpstreamType?: string;
	
	    static createFrom(source: any = {}) {
	        return new CheckParams(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ProxyList = source["ProxyList"];
	        this.ProxyType = source["ProxyType"];
	        this.Endpoint = source["Endpoint"];
	        this.Threads = source["Threads"];
	        this.UpstreamProxy = source["UpstreamProxy"];
	        this.UpstreamType = source["UpstreamType"];
	    }
	}

}

