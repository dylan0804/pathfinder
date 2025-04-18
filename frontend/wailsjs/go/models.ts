export namespace main {
	
	export class FileResult {
	    id: number;
	    path: string;
	    name: string;
	    extension: string;
	    isDir: boolean;
	    size: number;
	
	    static createFrom(source: any = {}) {
	        return new FileResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.path = source["path"];
	        this.name = source["name"];
	        this.extension = source["extension"];
	        this.isDir = source["isDir"];
	        this.size = source["size"];
	    }
	}
	export class SearchOptions {
	    pattern: string;
	    directory: string;
	    includeHidden: boolean;
	
	    static createFrom(source: any = {}) {
	        return new SearchOptions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.pattern = source["pattern"];
	        this.directory = source["directory"];
	        this.includeHidden = source["includeHidden"];
	    }
	}

}

