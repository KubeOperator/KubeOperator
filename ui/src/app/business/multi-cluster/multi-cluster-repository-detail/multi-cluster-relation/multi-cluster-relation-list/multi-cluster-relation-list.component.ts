import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {MultiClusterRepository, RelatedCluster} from "../../../multi-cluster-repository";
import {MultiClusterRepositoryService} from "../../../multi-cluster-repository.service";
import {CommonAlertService} from "../../../../../layout/common-alert/common-alert.service";
import {AlertLevels} from "../../../../../layout/common-alert/alert";

@Component({
    selector: 'app-multi-cluster-relation-list',
    templateUrl: './multi-cluster-relation-list.component.html',
    styleUrls: ['./multi-cluster-relation-list.component.css']
})
export class MultiClusterRelationListComponent implements OnInit {

    loading = false;
    selected: RelatedCluster[] = [];
    items: RelatedCluster[] = [];
    @Input() currentRepository: MultiClusterRepository;
    @Output() createEvent = new EventEmitter();
    @Output() deleteEvent = new EventEmitter();

    constructor(private multiClusterRepositoryService: MultiClusterRepositoryService, private commonAlertService: CommonAlertService) {
    }

    ngOnInit(): void {
        this.listRelationClusters();
    }

    onCreate() {
        this.createEvent.emit();

    }

    onDelete() {
        const names = [];
        for (const cluster of this.items) {
            names.push(cluster.clusterName);
        }
        this.deleteEvent.emit(names);
    }

    refresh() {
        this.listRelationClusters();
    }

    listRelationClusters() {
        this.loading = true;
        this.multiClusterRepositoryService.listRelations(this.currentRepository.name).subscribe(data => {
            this.loading = false;
            this.items = data;
        }, error => {
            this.loading = false;
            this.commonAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }

}
