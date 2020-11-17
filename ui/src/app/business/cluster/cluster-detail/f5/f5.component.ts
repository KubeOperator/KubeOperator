import {Component, OnInit} from '@angular/core';
import {F5Service} from './f5.service';
import {ActivatedRoute} from '@angular/router';
import {HttpClient} from '@angular/common/http';
import {Cluster} from '../../cluster';
import {F5} from './f5';
import {CommonAlertService} from '../../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {AlertLevels} from '../../../../layout/common-alert/alert';


@Component({
    selector: 'app-f5',
    templateUrl: './f5.component.html',
    styleUrls: ['./f5.component.css']
})
export class F5Component implements OnInit {
    item: F5 = new F5();
    createItem: F5 = new F5();
    currentCluster: Cluster;
    loading = false;

    constructor(
        private f5Service: F5Service,
        private route: ActivatedRoute,
        private  http: HttpClient,
        private commonAlertService: CommonAlertService,
        private translateService: TranslateService,

    ) {
    }

    ngOnInit(): void {
        this.route.parent.data.subscribe(data => {
          this.currentCluster = data.cluster;
          this.f5Service.getItems(this.currentCluster.name).subscribe(d => {
                this.item = d;
                this.item.clusterName = this.currentCluster.name;
                this.loading = true;
            });
        });
    }

    onSubmit() {
        this.loading = false;
        this.createItem = this.item;
        this.f5Service.create(this.createItem).subscribe(d => {
            if ( d.status === 'Running' ) {
                this.commonAlertService.showAlert(this.translateService.instant('APP_ADD_SUCCESS'), AlertLevels.SUCCESS);
                this.ngOnInit();
            }
        }, error => {
                this.commonAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
            });
    }

    onUpdate() {
            this.createItem = this.item;
            this.f5Service.update(this.createItem).subscribe(d => {
                if ( d.status === 'Running' ) {
                    this.commonAlertService.showAlert(this.translateService.instant('APP_UPDATE_SUCCESS'), AlertLevels.SUCCESS);
                    this.loading = true;
                }
            }, error => {
                    this.commonAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
                });
    }
}
