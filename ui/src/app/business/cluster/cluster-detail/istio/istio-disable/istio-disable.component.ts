import {Component, EventEmitter, Input, OnInit, Output, ViewChild} from '@angular/core';
import {IstioHelper} from '../istios';
import {Cluster} from '../../../cluster';
import {IstioService} from '../istio.service';
import {TranslateService} from '@ngx-translate/core';
import {AlertLevels} from '../../../../../layout/common-alert/alert';
import {CommonAlertService} from '../../../../../layout/common-alert/common-alert.service';

@Component({
    selector: 'app-istio-disable',
    templateUrl: './istio-disable.component.html',
    styleUrls: ['./istio-disable.component.css']
})
export class IstioDisableComponent implements OnInit {

    constructor(private istioService: IstioService,
                private translateService: TranslateService,
                private commonAlertService: CommonAlertService,
    ) {}
    opened = false;
    isSubmitGoing = false;
    items: IstioHelper[] = [];
    @Input() currentCluster: Cluster;
    @Output() disabled = new EventEmitter();


    ngOnInit(): void {
    }

    onSubmit() {
        this.isSubmitGoing = true;
        this.istioService.disable(this.currentCluster.name, this.items).subscribe(data => {
            this.commonAlertService.showAlert(this.translateService.instant('APP_ISTIO_STOP_SUCCESS'), AlertLevels.SUCCESS);
            this.opened = false;
            this.isSubmitGoing = false;
            this.disabled.emit();
        }, error => {
            this.isSubmitGoing = false;
        });
    }
    
    onCancel() {
        this.opened = false;
    }

    open(items: IstioHelper[]) {
        this.opened = true;
        this.items = items;
    }
}
