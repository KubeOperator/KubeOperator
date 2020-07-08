import {Component, OnInit} from '@angular/core';
import {ActivatedRoute} from "@angular/router";
import {Cluster} from "../../../cluster";
import {ToolsService} from "../../tools/tools.service";

@Component({
    selector: 'app-registry',
    templateUrl: './registry.component.html',
    styleUrls: ['./registry.component.css']
})
export class RegistryComponent implements OnInit {

    currentCluster: Cluster;
    ready = false;
    toolName = 'registry';

    constructor(private toolService: ToolsService, private route: ActivatedRoute) {
    }

    ngOnInit(): void {
        this.route.parent.parent.data.subscribe(data => {
            this.currentCluster = data.cluster;
            this.toolService.list(this.currentCluster.name).subscribe(d => {
                for (const tool of d) {
                    if (tool.name === this.toolName) {
                        this.ready = tool.status === 'Running';
                    }
                }
            });
        });
    }

}
