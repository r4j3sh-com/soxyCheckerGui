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

export namespace config {
	
	export class Config {
	    lastProxyType: string;
	    lastEndpoint: string;
	    lastThreadCount: number;
	    lastUpstreamProxy: string;
	    lastUpstreamProxyType: string;
	    defaultEndpoints: string[];
	    maxThreads: number;
	    theme: string;
	    enableGeolocation: boolean;
	    exportFormat: string;
	    autoSaveResults: boolean;
	    autoSavePath: string;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.lastProxyType = source["lastProxyType"];
	        this.lastEndpoint = source["lastEndpoint"];
	        this.lastThreadCount = source["lastThreadCount"];
	        this.lastUpstreamProxy = source["lastUpstreamProxy"];
	        this.lastUpstreamProxyType = source["lastUpstreamProxyType"];
	        this.defaultEndpoints = source["defaultEndpoints"];
	        this.maxThreads = source["maxThreads"];
	        this.theme = source["theme"];
	        this.enableGeolocation = source["enableGeolocation"];
	        this.exportFormat = source["exportFormat"];
	        this.autoSaveResults = source["autoSaveResults"];
	        this.autoSavePath = source["autoSavePath"];
	    }
	}

}

