import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { StorageProvisionerSyncComponent } from './storage-provisioner-sync.component';

describe('StorageProvisionerSyncComponent', () => {
  let component: StorageProvisionerSyncComponent;
  let fixture: ComponentFixture<StorageProvisionerSyncComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ StorageProvisionerSyncComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(StorageProvisionerSyncComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
