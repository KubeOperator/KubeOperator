import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {SecurityService} from "../security.service";
import {Cluster} from "../../../cluster";
import {CommonAlertService} from "../../../../../layout/common-alert/common-alert.service";
import {AlertLevels} from "../../../../../layout/common-alert/alert";

@Component({
    selector: 'app-security-task-create',
    templateUrl: './security-task-create.component.html',
    styleUrls: ['./security-task-create.component.css']
})
export class SecurityTaskCreateComponent implements OnInit {

    constructor(private service: SecurityService, private alertService: CommonAlertService) {
    }

    opened = false;
    isSubmitGoing = false;
    @Input() currentCluster: Cluster;
    @Output() created = new EventEmitter();

    ngOnInit(): void {
    }

    onCancel() {
        this.opened = false;
    }

    onSubmit() {
        this.service.create(this.currentCluster.name).subscribe(data => {
            this.created.emit();
            this.opened = false;
        }, error => {
            this.alertService.showAlert(error.error.msg, AlertLevels.ERROR);
            this.opened = false;
        });
    }

    open() {
        this.opened = true;
    }

}
