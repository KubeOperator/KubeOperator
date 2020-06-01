import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {HostService} from '../host.service';
import {BaseModelComponent} from '../../../shared/class/BaseModelComponent';
import {Host} from '../host';

@Component({
    selector: 'app-host-list',
    templateUrl: './host-list.component.html',
    styleUrls: ['./host-list.component.css']
})
export class HostListComponent extends BaseModelComponent<Host> implements OnInit {

    @Output() detailEvent = new EventEmitter<Host>();

    constructor(private hostService: HostService) {
        super(hostService);
    }

    ngOnInit(): void {
        super.ngOnInit();
    }

    onStatusDetail(item) {
        this.detailEvent.emit(item);
    }

}
