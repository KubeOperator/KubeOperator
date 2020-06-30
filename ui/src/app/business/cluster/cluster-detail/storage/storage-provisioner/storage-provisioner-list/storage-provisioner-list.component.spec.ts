import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { StorageProvisionerListComponent } from './storage-provisioner-list.component';

describe('StorageProvisionerListComponent', () => {
  let component: StorageProvisionerListComponent;
  let fixture: ComponentFixture<StorageProvisionerListComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ StorageProvisionerListComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(StorageProvisionerListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
