import {Component, OnInit} from '@angular/core';
import {BaseModelDirective} from '../../../../shared/class/BaseModelDirective';
import {VmConfig} from '../vm-config';
import {VmConfigService} from '../vm-config.service';

@Component({
    selector: 'app-vm-config-list',
    templateUrl: './vm-config-list.component.html',
    styleUrls: ['./vm-config-list.component.css']
})
export class VmConfigListComponent extends BaseModelDirective<VmConfig> implements OnInit {

    loading = false;

    constructor(private vmConfigService: VmConfigService) {
        super(vmConfigService);
    }

    ngOnInit(): void {
        super.ngOnInit();
    }

}
