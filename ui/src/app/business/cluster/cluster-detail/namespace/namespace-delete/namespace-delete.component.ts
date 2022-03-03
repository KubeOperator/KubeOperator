import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {KubernetesService} from '../../../kubernetes.service';
import {ActivatedRoute} from '@angular/router';
import {CommonAlertService} from '../../../../../layout/common-alert/common-alert.service';
import {ModalAlertService} from '../../../../../shared/common-component/modal-alert/modal-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {Cluster} from '../../../cluster';
import {AlertLevels} from '../../../../../layout/app-alert/alert';
import {ToolsService} from '../../tools/tools.service';

@Component({
    selector: 'app-namespace-delete',
    templateUrl: './namespace-delete.component.html',
    styleUrls: ['./namespace-delete.component.css']
})
export class NamespaceDeleteComponent implements OnInit {


    opened = false;
    namespace;
    @Output() deleted = new EventEmitter();
    @Input() currentCluster: Cluster;

    constructor(private service: KubernetesService, private route: ActivatedRoute,
                private commonAlertService: CommonAlertService,
                private modalAlertService: ModalAlertService,
                private translateService: TranslateService,
                private toolsService: ToolsService) {
    }

    ngOnInit(): void {
    }


    open(namespace) {
        this.namespace = namespace;
        this.opened = true;
    }

    onSubmit() {
        let exitStr: string = '';
        this.toolsService.list(this.currentCluster.name).subscribe(d => {
            if (d) {
                for (const tool of d) {
                    if (tool.vars !== null) {
                        if (tool.vars["namespace"] === this.namespace && tool.status !== 'Waiting') {
                            exitStr += tool.name + ',';
                        }
                    }
                }
                if (exitStr === '') {
                    this.service.deleteNamespace(this.currentCluster.name, this.namespace).subscribe(res => {
                        this.opened = false;
                        this.commonAlertService.showAlert(this.translateService.instant('APP_DELETE_SUCCESS'), AlertLevels.SUCCESS);
                        this.deleted.emit();
                    }, error => {
                        this.opened = false;
                        this.commonAlertService.showAlert(error.error.message, AlertLevels.ERROR);
                    });
                } else {
                    exitStr = exitStr.substring(0, exitStr.length - 1);
                    this.modalAlertService.showAlert(this.translateService.instant('APP_DELETE_NS_ERR') + exitStr, AlertLevels.ERROR);
                }
            }
        });
    }

    onCancel() {
        this.opened = false;
    }
}

