import {Component, EventEmitter, Input, OnInit, Output, ViewChild} from '@angular/core';
import {MultiClusterRepository, RelateClusterRequest} from "../../../multi-cluster-repository";
import {ClusterService} from "../../../../cluster/cluster.service";
import {MultiClusterRepositoryService} from "../../../multi-cluster-repository.service";
import {ModalAlertService} from "../../../../../shared/common-component/modal-alert/modal-alert.service";
import {AlertLevels} from "../../../../../layout/common-alert/alert";
import {NgForm} from "@angular/forms";

@Component({
    selector: 'app-multi-cluster-relation-create',
    templateUrl: './multi-cluster-relation-create.component.html',
    styleUrls: ['./multi-cluster-relation-create.component.css']
})
export class MultiClusterRelationCreateComponent implements OnInit {

    constructor(private clusterService: ClusterService,
                private multiClusterRepositoryService: MultiClusterRepositoryService, private modalAlertService: ModalAlertService) {
    }

    opened = false;
    isSubmitGoing = false;
    selections: any[] = [];
    clusters: any[] = [];
    @Input() currentRepository: MultiClusterRepository;
    @Output() created = new EventEmitter();
    @ViewChild('itemForm') itemForm: NgForm;
    options: any = {
        multiple: true,
    };

    ngOnInit(): void {
    }

    open() {
        this.listClusters();
        this.isSubmitGoing = false;
        this.selections = [];
        this.clusters = [];
        this.itemForm.resetForm();
        this.opened = true;
    }

    onCancel() {
        this.opened = false;
    }

    onSubmit() {
        if (this.isSubmitGoing) {
            return;
        }
        this.isSubmitGoing = true;
        const names = [];
        for (const c of this.clusters) {
            names.push(c.id);
        }
        this.multiClusterRepositoryService.createRelations(this.currentRepository.name, names).subscribe(data => {
            this.isSubmitGoing = false;
            this.opened = false;
            this.created.emit();
        }, error => {
            this.modalAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }

    listClusters() {
        this.clusterService.list().subscribe(data => {
            const s = [];
            for (const d of data.items.filter((item) => {
                return !item.multiClusterRepository && item.status === 'Running';
            })) {
                s.push({id: d.name, text: d.name});
            }
            this.selections = s;
        });
    }

}
