import {Component, OnInit} from '@angular/core';
import {Cluster} from "../../cluster";
import {ClusterTool} from "../tools/tools";
import {ToolsService} from "../tools/tools.service";
import {ActivatedRoute} from "@angular/router";

@Component({
    selector: 'app-catelog',
    templateUrl: './catalog.component.html',
    styleUrls: ['./catalog.component.css']
})
export class CatalogComponent implements OnInit {

    currentCluster: Cluster;
    ready = false;
    toolName = 'kubeapps';
    item: ClusterTool;

    constructor(private toolService: ToolsService, private route: ActivatedRoute) {
    }

    ngOnInit(): void {
        this.route.parent.data.subscribe(data => {
            this.currentCluster = data.cluster;
            this.toolService.list(this.currentCluster.name).subscribe(d => {
                for (const tool of d) {
                    if (tool.name === this.toolName) {
                        this.item = tool;
                        this.ready = tool.status === 'Running';
                    }
                }
            });
        });
    }
}
