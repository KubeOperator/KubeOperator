import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {MultiClusterRepository, MultiClusterSyncLog} from "../../../multi-cluster-repository";
import {MultiClusterRepositoryService} from "../../../multi-cluster-repository.service";

@Component({
    selector: 'app-multi-cluster-log-list',
    templateUrl: './multi-cluster-log-list.component.html',
    styleUrls: ['./multi-cluster-log-list.component.css']
})
export class MultiClusterLogListComponent implements OnInit {

    loading = false;
    items: MultiClusterSyncLog[] = [];
    page = 1;
    size = 10;
    total = 0;
    @Input() currentRepository: MultiClusterRepository;
    @Output() detailEvent = new EventEmitter<string>();


    constructor(private multiClusterRepositoryService: MultiClusterRepositoryService) {
    }

    ngOnInit(): void {
        this.refresh();
    }

    refresh() {
        this.loading = true;
        this.multiClusterRepositoryService.getLog(this.currentRepository.name, this.page, this.size).subscribe(data => {
            this.items = data.items;
            this.total = data.total;
            this.loading = false;
        });
    }

    onDetail(item: MultiClusterSyncLog) {
        this.detailEvent.emit(item.id);
    }
}
