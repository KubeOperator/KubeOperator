import {Component, OnInit, ViewChild} from '@angular/core';
import {MultiClusterRepository, MultiClusterRepositoryUpdateRequest} from "../../multi-cluster-repository";
import {ActivatedRoute} from "@angular/router";
import {MultiClusterRepositoryService} from "../../multi-cluster-repository.service";
import {NgForm} from "@angular/forms";
import {CommonAlertService} from "../../../../layout/common-alert/common-alert.service";
import {AlertLevels} from "../../../../layout/common-alert/alert";

@Component({
    selector: 'app-multi-cluster-setting',
    templateUrl: './multi-cluster-setting.component.html',
    styleUrls: ['./multi-cluster-setting.component.css']
})
export class MultiClusterSettingComponent implements OnInit {

    constructor(private route: ActivatedRoute, private multiClusterRepositoryService: MultiClusterRepositoryService,
                private commonAlertService: CommonAlertService) {
    }

    currentRepository: MultiClusterRepository;
    item: MultiClusterRepositoryUpdateRequest = new MultiClusterRepositoryUpdateRequest();
    @ViewChild('itemForm')
    itemForm: NgForm;
    isSubmitGoing = false;

    ngOnInit(): void {
        this.route.parent.data.subscribe(d => {
            this.currentRepository = d.repo;
            this.refresh();
        });
    }

    onReset() {
        this.itemForm.resetForm();
        this.refresh();
    }

    onSubmit() {
        if (this.isSubmitGoing) {
            return;
        }
        this.isSubmitGoing = true;
        this.multiClusterRepositoryService.update(this.currentRepository.name, this.item).subscribe(data => {
            this.refresh();
            this.isSubmitGoing = false;
        }, error => {
            this.commonAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
            this.isSubmitGoing = false;
        });
    }


    refresh() {
        this.multiClusterRepositoryService.get(this.currentRepository.name).subscribe(data => {
            this.item.gitTimeout = data.gitTimeout;
            this.item.syncInterval = data.syncInterval;
            this.item.syncEnable = data.syncEnable;
        }, error => {
            this.commonAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }

}
