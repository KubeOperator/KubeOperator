import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {BaseModelComponent} from "../../../../shared/class/BaseModelComponent";
import {Region, RegionCreateRequest} from "../region";
import {RegionService} from "../region.service";
import {ModalAlertService} from "../../../../shared/common-component/modal-alert/modal-alert.service";
import {CommonAlertService} from "../../../../layout/common-alert/common-alert.service";
import {NgForm} from "@angular/forms";
import {CloudProviderService} from "../cloud-provider.service";
import {CloudProvider} from "../cloud-provider";
import {AlertLevels} from "../../../../layout/common-alert/alert";

@Component({
    selector: 'app-region-create',
    templateUrl: './region-create.component.html',
    styleUrls: ['./region-create.component.css']
})
export class RegionCreateComponent extends BaseModelComponent<Region> implements OnInit {

    opened = false;
    isSubmitGoing = false;
    item: RegionCreateRequest = new RegionCreateRequest();
    cloudProviders: CloudProvider[] = [];
    @Output() created = new EventEmitter();
    @ViewChild('regionForm', {static: true}) regionForm: NgForm;
    @ViewChild('paramsForm', {static: true}) paramsForm: NgForm;


    constructor(private regionService: RegionService, private modalAlertService: ModalAlertService,
                private commonAlertService: CommonAlertService, private cloudProviderService: CloudProviderService) {
        super(regionService);
    }

    ngOnInit(): void {
    }

    open() {
        this.cloudProviderService.list().subscribe(res => {
            this.cloudProviders = res.items;
            this.opened = true;
            this.item = new RegionCreateRequest();
        }, error => {
            this.modalAlertService.showAlert("", AlertLevels.ERROR);
        })

    }

    onCancel() {
        this.opened = false;
    }

    onCheckParams() {
        this.item.regionVars['provider'] = this.item.cloudProvider
        this.regionService.checkValid(this.item).subscribe(data => {
            console.log("success")
        }, error => {
            console.log("failed")
        })
    }
}
