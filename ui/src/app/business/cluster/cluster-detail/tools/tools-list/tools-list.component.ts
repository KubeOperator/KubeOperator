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
            let logIndex = -1;
            let lokiIndex = -1;
            let prometheusIndex = -1;
            let grafanaIndex = -1;
            for (let i = 0; i < data.length; i++) {
                data[i].isDisable = false;
                switch (data[i].name) {
                    case 'logging': 
                        logIndex = i;
                        break;
                    case 'loki': 
                        lokiIndex = i;
                        break;
                    case 'prometheus': 
                        prometheusIndex = i;
                        break;
                    case 'grafana': 
                        grafanaIndex = i;
                        break;
                }
            }
            if (logIndex !== -1 && lokiIndex !== -1) {
                data[logIndex].isDisable = (data[lokiIndex].status !== 'Waiting');
                data[lokiIndex].isDisable = (data[logIndex].status !== 'Waiting');
            }
            if (prometheusIndex !== -1 && grafanaIndex !== -1) {
                data[grafanaIndex].isDisable = (data[prometheusIndex].status !== 'Running');
            } else if (grafanaIndex !== -1) {
                data[grafanaIndex].isDisable = true;
            }
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
