import {Component, OnInit} from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import {Cluster} from '../../cluster';
import {ToolsService} from "../tools/tools.service";
import {ClusterTool} from "../tools/tools";

@Component({
    selector: 'app-monitor',
    templateUrl: './monitor.component.html',
    styleUrls: ['./monitor.component.css']
})
export class MonitorComponent implements OnInit {

    currentCluster: Cluster;
    ready = false;
    toolName = 'prometheus';
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
