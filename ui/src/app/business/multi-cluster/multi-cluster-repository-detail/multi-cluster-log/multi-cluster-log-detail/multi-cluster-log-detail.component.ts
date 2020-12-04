import {Component, Input, OnInit} from '@angular/core';
import {MultiClusterRepository, MultiClusterSyncLogDetail} from "../../../multi-cluster-repository";
import {MultiClusterRepositoryService} from "../../../multi-cluster-repository.service";

@Component({
    selector: 'app-multi-cluster-log-detail',
    templateUrl: './multi-cluster-log-detail.component.html',
    styleUrls: ['./multi-cluster-log-detail.component.css']
})
export class MultiClusterLogDetailComponent implements OnInit {

    constructor(private multiClusterRepositoryService: MultiClusterRepositoryService) {
    }

    opened = false;
    item: MultiClusterSyncLogDetail = new MultiClusterSyncLogDetail();
    logId: string;
    @Input() currentRepository: MultiClusterRepository;

    ngOnInit(): void {
    }

    open(logId: string) {
        this.opened = true;
        this.logId = logId;
        this.refresh();
    }

    refresh() {
        this.multiClusterRepositoryService.getLogDetail(this.currentRepository.name, this.logId).subscribe(data => {
            console.log(data);
            this.item = data;
        });
    }
}
