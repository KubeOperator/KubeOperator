import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { HostSyncComponent } from './host-sync.component';

describe('HostSyncComponent', () => {
  let component: HostSyncComponent;
  let fixture: ComponentFixture<HostSyncComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ HostSyncComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(HostSyncComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
