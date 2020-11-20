import {Component, EventEmitter, Input, OnDestroy, OnInit, Output} from '@angular/core';
import {ClusterTool} from '../tools';
import {ToolsService} from '../tools.service';
import {Cluster} from '../../../cluster';

@Component({
    selector: 'app-tools-list',
    templateUrl: './tools-list.component.html',
    styleUrls: ['./tools-list.component.css']
})
export class ToolsListComponent implements OnInit, OnDestroy {

    constructor(private service: ToolsService) {
    }


    items: ClusterTool[] = [];
    timer;
    @Input() currentCluster: Cluster;
    @Output() enableEvent = new EventEmitter<ClusterTool>();
    @Output() disableEvent = new EventEmitter<ClusterTool>();
    @Output() failedEvent = new EventEmitter<ClusterTool>();

    ngOnInit(): void {
        this.refresh();
    }

    ngOnDestroy(): void {
        clearInterval(this.timer);
    }

    refresh() {
        this.service.list(this.currentCluster.name).subscribe(data => {
            this.items = data;
        });
    }

    onEnable(item: ClusterTool) {
        this.enableEvent.emit(item);
    }

    onDisable(item: ClusterTool) {
        this.disableEvent.emit(item);
    }

    onFailed(item: ClusterTool) {
        this.failedEvent.emit(item);
    }

    openFrame(item: ClusterTool) {
        window.open(item.url.replace('{cluster_name}', this.currentCluster.name), '_blank');
    }

    polling() {
        this.timer = setInterval(() => {
            this.refresh();
        }, 5000);
    }

}
