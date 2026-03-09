export namespace main {
	
	export class AppConfig {
	    frpcPath: string;
	    localPorts: number[];
	    localPort?: number;
	    serverAddr: string;
	    serverPort: number;
	    authToken: string;
	
	    static createFrom(source: any = {}) {
	        return new AppConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.frpcPath = source["frpcPath"];
	        this.localPorts = source["localPorts"];
	        this.localPort = source["localPort"];
	        this.serverAddr = source["serverAddr"];
	        this.serverPort = source["serverPort"];
	        this.authToken = source["authToken"];
	    }
	}
	export class RuntimeState {
	    running: boolean;
	    pid: number;
	    appDir: string;
	    configFile: string;
	    frpcToml: string;
	    frpcPath: string;
	    lastError: string;
	
	    static createFrom(source: any = {}) {
	        return new RuntimeState(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.running = source["running"];
	        this.pid = source["pid"];
	        this.appDir = source["appDir"];
	        this.configFile = source["configFile"];
	        this.frpcToml = source["frpcToml"];
	        this.frpcPath = source["frpcPath"];
	        this.lastError = source["lastError"];
	    }
	}

}

