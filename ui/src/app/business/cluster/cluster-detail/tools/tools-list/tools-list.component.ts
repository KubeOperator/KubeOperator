import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
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

    constructor(private service: ToolsService) {
    }

    items: ClusterTool[] = [];
    @Input() currentCluster: Cluster;
    @Output() enableEvent = new EventEmitter<ClusterTool>();

    ngOnInit(): void {
        this.refresh();
    }


    refresh() {
        this.service.list(this.currentCluster.name).subscribe(data => {
            this.items = data;
        });
    }

    onEnable(item: ClusterTool) {
        this.enableEvent.emit(item);
    }

}
