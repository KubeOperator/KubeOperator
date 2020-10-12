import { ComponentFixture, TestBed } from '@angular/core/testing';

import { VmConfigListComponent } from './vm-config-list.component';

describe('VmConfigListComponent', () => {
  let component: VmConfigListComponent;
  let fixture: ComponentFixture<VmConfigListComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ VmConfigListComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(VmConfigListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
