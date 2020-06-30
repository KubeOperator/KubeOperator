import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { StorageProvisionerComponent } from './storage-provisioner.component';

describe('StorageProvisionerComponent', () => {
  let component: StorageProvisionerComponent;
  let fixture: ComponentFixture<StorageProvisionerComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ StorageProvisionerComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(StorageProvisionerComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
