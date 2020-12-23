import {Component, OnInit} from '@angular/core';
import {IstioService} from "./istio.service";
import {ActivatedRoute} from "@angular/router";
import {Cluster} from "../../cluster";
import {IstioHelper} from './istios';
import {AlertLevels} from '../../../../layout/common-alert/alert';
import {CommonAlertService} from '../../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';

@Component({
    selector: 'app-istios',
    templateUrl: './istio.component.html',
    styleUrls: ['./istio.component.css']
})
export class IstioComponent implements OnInit {

    constructor(
        private istioService: IstioService,
        private route: ActivatedRoute,
        private commonAlertService: CommonAlertService,
        private translateService: TranslateService
    ) {}

    accordionLoading: boolean = false;
    enableLoading: boolean = false;
    disableLoading: boolean = false;
    currentCluster: Cluster;
    stepOpen: boolean = true;
    baseCfg: IstioHelper = new IstioHelper;
    pilotCfg: IstioHelper = new IstioHelper;
    ingressCfg: IstioHelper = new IstioHelper;
    egressCfg: IstioHelper = new IstioHelper;

    ngOnInit(): void {
        this.accordionLoading = true;
        this.route.parent.data.subscribe(data => {
            this.currentCluster = data.cluster;
        });
        this.refresh();
        this.accordionLoading = false;
    }
    submit (operation: string) {
        this.enableLoading = true;
        var items: IstioHelper[] = [];
        if (operation === 'start') {
            this.baseCfg.enable = true;
            this.getOperation(items, this.baseCfg);
            this.pilotCfg.enable = true;
            this.getOperation(items, this.pilotCfg);
        }
        this.getOperation(items, this.ingressCfg);
        this.getOperation(items, this.egressCfg);
        this.istioService.enable(this.currentCluster.name, items).subscribe(data => {
            this.commonAlertService.showAlert(this.translateService.instant('APP_RESTORE_START_SUCCESS'), AlertLevels.SUCCESS);
            this.enableLoading = false;
        }, error => {
            this.enableLoading = false;
        });
    }
    stopIstio () {
        this.disableLoading = true;
        var items: IstioHelper[] = [];
        this.disAble(items, this.baseCfg);
        this.disAble(items, this.pilotCfg);
        this.disAble(items, this.ingressCfg);
        this.disAble(items, this.egressCfg);
        this.istioService.disable(this.currentCluster.name, items).subscribe(data => {
            this.commonAlertService.showAlert(this.translateService.instant('APP_RESTORE_START_SUCCESS'), AlertLevels.SUCCESS);
            this.disableLoading = false;
        }, error => {
            this.disableLoading = false;
        });
    }
    disAble(items: IstioHelper[], istio: IstioHelper) {
        istio.operation = 'disable';
        istio.enable = false;
        items.push(istio);
    }
    getOperation(items: IstioHelper[], istio: IstioHelper) {
        if (istio.enable && istio.cluster_istio.status !== 'Running') {
            istio.operation = 'enable';
            items.push(istio);
        } else if (!istio.enable && istio.cluster_istio.status !== 'Waiting') {
            istio.operation = 'disable';
            items.push(istio);
        }
    }
    refresh() {
        this.istioService.list(this.currentCluster.name).subscribe(data => {
            for (const item of data) {
                switch (item.cluster_istio.name) {
                    case 'base':
                        this.baseCfg.cluster_istio = item.cluster_istio;
                        this.baseCfg.enable = (item.cluster_istio.status === "Running");
                        break;
                    case 'pilot':
                        this.pilotCfg.cluster_istio = item.cluster_istio;
                        this.pilotCfg.enable = (item.cluster_istio.status === "Running");
                        break;
                    case 'ingress':
                        this.ingressCfg.cluster_istio = item.cluster_istio;
                        this.ingressCfg.enable = (item.cluster_istio.status === "Running");
                        break;
                    case 'egress':
                        this.egressCfg.cluster_istio = item.cluster_istio;
                        this.egressCfg.enable = (item.cluster_istio.status === "Running");
                        break;
                }
            }
        });
        this.setDefaultVars();
    }
    setDefaultVars() {
        this.baseCfg.vars = {
            'global.istiod.enableAnalysis': true,
        }
        this.pilotCfg.vars = {
            'pilot.resources.requests.cpu': 500,
            'pilot.resources.requests.memory': 2048, 
            'pilot.resources.limits.cpu': 1000,
            'pilot.resources.limits.memory': 4096,
            'pilot.traceSampling': 1,
        };
        this.ingressCfg.vars = {
            'gateways.istio-ingressgateway.type': 'NodePort',
            'gateways.istio-ingressgateway.resources.requests.cpu': 100,
            'gateways.istio-ingressgateway.resources.requests.memory': 128,
            'gateways.istio-ingressgateway.resources.limits.cpu': 2000,
            'gateways.istio-ingressgateway.resources.limits.memory': 1024,
        };
        this.egressCfg.vars = {
            'gateways.istio-egressgateway.type': 'NodePort',
            'gateways.istio-egressgateway.resources.requests.cpu': 100,
            'gateways.istio-egressgateway.resources.requests.memory': 128,
            'gateways.istio-egressgateway.resources.limits.cpu': 2000,
            'gateways.istio-egressgateway.resources.limits.memory': 1024,
        };
    }
}
