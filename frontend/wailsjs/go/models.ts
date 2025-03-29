export namespace main {
	
	export class FileResult {
	    path: string;
	    name: string;
	    directory: string;
	    isDir: boolean;
	    size: number;
	    // Go type: time
	    modTime: any;
	
	    static createFrom(source: any = {}) {
	        return new FileResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.name = source["name"];
	        this.directory = source["directory"];
	        this.isDir = source["isDir"];
	        this.size = source["size"];
	        this.modTime = this.convertValues(source["modTime"], null);
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
	export class SearchOptions {
	    pattern: string;
	    directory: string;
	    caseSensitive: boolean;
	    includeHidden: boolean;
	    fileTypes: string[];
	    maxDepth: number;
	    extensions: string[];
	
	    static createFrom(source: any = {}) {
	        return new SearchOptions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.pattern = source["pattern"];
	        this.directory = source["directory"];
	        this.caseSensitive = source["caseSensitive"];
	        this.includeHidden = source["includeHidden"];
	        this.fileTypes = source["fileTypes"];
	        this.maxDepth = source["maxDepth"];
	        this.extensions = source["extensions"];
	    }
	}

}

