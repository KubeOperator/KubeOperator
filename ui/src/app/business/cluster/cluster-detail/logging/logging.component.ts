import {Component, OnInit} from '@angular/core';
import {ToolsService} from "../tools/tools.service";
import {ActivatedRoute} from "@angular/router";
import {Cluster} from "../../cluster";

@Component({
    selector: 'app-logging',
    templateUrl: './logging.component.html',
    styleUrls: ['./logging.component.css']
})
export class LoggingComponent implements OnInit {

    constructor(private toolService: ToolsService, private route: ActivatedRoute) {
    }

    currentCluster: Cluster;
    ready = false;
    toolName = 'EFK';

    ngOnInit(): void {
        this.route.parent.data.subscribe(data => {
            this.currentCluster = data.cluster;
            this.toolService.list(this.currentCluster.name).subscribe(d => {
                for (const tool of d) {
                    if (tool.name === this.toolName) {
                        this.ready = tool.status === 'running';
                    }
                }
            });
        });
    }

}
