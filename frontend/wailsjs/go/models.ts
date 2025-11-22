export namespace main {
	
	export class VerbConfig {
	    value: string;
	    enabled: boolean;
	
	    static createFrom(source: any = {}) {
	        return new VerbConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.value = source["value"];
	        this.enabled = source["enabled"];
	    }
	}
	export class Config {
	    inputPath: string;
	    inputMode: string;
	    inputFormat: string;
	    ragged: boolean;
	    headerless: boolean;
	    fieldSeparator: string;
	    outputFormat: string;
	    verbs: VerbConfig[];
	    options: string;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.inputPath = source["inputPath"];
	        this.inputMode = source["inputMode"];
	        this.inputFormat = source["inputFormat"];
	        this.ragged = source["ragged"];
	        this.headerless = source["headerless"];
	        this.fieldSeparator = source["fieldSeparator"];
	        this.outputFormat = source["outputFormat"];
	        this.verbs = this.convertValues(source["verbs"], VerbConfig);
	        this.options = source["options"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

