import {Component, OnInit} from '@angular/core';
import {ClusterTool} from "../tools";
import {ToolsService} from "../tools.service";
import {Cluster} from "../../../cluster";
import {ActivatedRoute} from "@angular/router";

@Component({
    selector: 'app-tools-list',
    templateUrl: './tools-list.component.html',
    styleUrls: ['./tools-list.component.css']
})
export class ToolsListComponent implements OnInit {

    constructor(private service: ToolsService, private route: ActivatedRoute) {
    }

    items: ClusterTool[] = [];
    currentCluster: Cluster;

    ngOnInit(): void {
        this.route.parent.data.subscribe(data => {
            this.currentCluster = data.cluster;
            this.refresh();
        });
    }


    refresh() {
        this.service.list(this.currentCluster.name).subscribe(data => {
            this.items = data;
        });
    }

}
