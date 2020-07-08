import {Component, OnInit} from '@angular/core';
import {Cluster} from "../../../cluster";
import {ActivatedRoute} from "@angular/router";
import {ToolsService} from "../../tools/tools.service";

@Component({
    selector: 'app-chartmuseum',
    templateUrl: './chartmuseum.component.html',
    styleUrls: ['./chartmuseum.component.css']
})
export class ChartmuseumComponent implements OnInit {

    currentCluster: Cluster;
    ready = false;
    toolName = 'chartmuseum';

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
